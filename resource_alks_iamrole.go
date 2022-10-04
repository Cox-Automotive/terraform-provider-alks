package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func resourceAlksIamRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlksIamRoleCreate,
		ReadContext:   resourceAlksIamRoleRead,
		UpdateContext: resourceAlksIamRoleUpdate,
		DeleteContext: resourceAlksIamRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		MigrateState:  migrateState,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name_prefix"},
				ValidateFunc:  ValidRoleName,
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
				ValidateFunc:  ValidRolePrefix,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"assume_role_policy", "type"},
			},
			"assume_role_policy": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ExactlyOneOf:     []string{"assume_role_policy", "type"},
				DiffSuppressFunc: SuppressEquivalentTrustPolicyDiffs,
			},
			"include_default_policies": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"role_added_to_ip": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_alks_access": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
			"template_fields": {
				Type:     schema.TypeMap,
				Elem:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"max_session_duration_in_seconds": {
				Type:     schema.TypeInt,
				Default:  3600,
				Optional: true,
				ForceNew: true,
			},
			"tags":     TagsSchema(),
			"tags_all": TagsSchemaComputed(),
		},
		CustomizeDiff: customdiff.All(
			SetTagsDiff,
			trustPoliciesWithIncludeDefaultPolicies,
		),
	}
}

func resourceAlksIamRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS IAM Role Create")
	var roleName = NameWithPrefix(d.Get("name").(string), d.Get("name_prefix").(string))
	var incDefPol = d.Get("include_default_policies").(bool)
	var enableAlksAccess = d.Get("enable_alks_access").(bool)
	var rawTemplateFields = d.Get("template_fields").(map[string]interface{})
	var maxSessionDurationInSeconds = d.Get("max_session_duration_in_seconds").(int)
	var tags = d.Get("tags").(map[string]interface{})

	templateFields := make(map[string]string)
	for k, v := range rawTemplateFields {
		templateFields[k] = v.(string)
	}

	include := false
	if incDefPol {
		include = true
	}

	providerStruct := meta.(*AlksClient)
	client := providerStruct.client

	//Role Specific tags will overwrite default tags if value is defined in both maps
	allTags := tagMapToSlice(combineTagMaps(providerStruct.defaultTags, tags))

	if err := validateIAMEnabled(client); err != nil {
		return diag.FromErr(err)
	}

	options := &alks.CreateIamRoleOptions{
		RoleName:                    &roleName,
		IncludeDefaultPolicies:      &include,
		AlksAccess:                  &enableAlksAccess,
		TemplateFields:              &templateFields,
		MaxSessionDurationInSeconds: &maxSessionDurationInSeconds,
		Tags:                        &allTags,
	}

	if roleType, ok := d.GetOk("type"); ok {
		roleTypeString := roleType.(string)
		options.RoleType = &roleTypeString
	} else {
		trustPolicyString := d.Get("assume_role_policy").(string)

		trustPolicy := new(map[string]interface{})

		json.Unmarshal([]byte(trustPolicyString), trustPolicy)

		options.TrustPolicy = trustPolicy
	}

	resp, err := client.CreateIamRole(options)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.RoleName)
	_ = d.Set("role_added_to_ip", resp.RoleAddedToIP)

	log.Printf("[INFO] alks_iamrole.id: %v", d.Id())

	return resourceAlksIamRoleRead(ctx, d, meta)
}

func resourceAlksIamRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS IAM Role Delete")

	providerStruct := meta.(*AlksClient)
	client := providerStruct.client
	if err := validateIAMEnabled(client); err != nil {
		return diag.FromErr(err)
	}

	if err := client.DeleteIamRole(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceAlksIamRoleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS IAM Role Read")
	providerStruct := meta.(*AlksClient)
	client := providerStruct.client

	defaultTags := providerStruct.defaultTags
	ignoreTags := providerStruct.ignoreTags

	// Check if role exists.
	if d.Id() == "" || d.Id() == "none" {
		return nil
	}

	foundRole, err := client.GetIamRole(d.Id())

	if err != nil {
		// If 404 Role not found error, an error and a role with Exists field set to false will come back from alks-go
		// We will log ther error and set id to "" and return nil, letting terraform decide how to handle role not found.
		if foundRole != nil {
			if foundRole.Exists != true {
				log.Printf("[Error] %#v", err)
				d.SetId("")
				return nil
			}
		}
		d.SetId("")
		return diag.FromErr(err)
	}

	log.Printf("[INFO] alks_iamrole.id %v", d.Id())

	_ = d.Set("name", foundRole.RoleName)
	_ = d.Set("name_prefix", NamePrefixFromName(foundRole.RoleName))
	_ = d.Set("arn", foundRole.RoleArn)
	_ = d.Set("ip_arn", foundRole.RoleIPArn)
	_ = d.Set("enable_alks_access", foundRole.AlksAccess)
	_ = d.Set("role_added_to_ip", foundRole.RoleAddedToIP)

	var roleType = d.Get("type")
	if (roleType == nil) || (roleType == "") {
		jsonStrPolicy, err := json.Marshal(foundRole.TrustPolicy)
		if err == nil {
			_ = d.Set("assume_role_policy", string(jsonStrPolicy))
		}
		_ = d.Set("include_default_policies", false)
	}
	_ = d.Set("max_session_duration_in_seconds", foundRole.MaxSessionDurationInSeconds)

	allTags := tagSliceToMap(foundRole.Tags)
	localTags := removeIgnoredTags(allTags, *ignoreTags)

	if err := d.Set("tags_all", localTags); err != nil {
		return diag.FromErr(err)
	}

	roleSpecificTags := removeDefaultTags(localTags, defaultTags)

	if err := d.Set("tags", roleSpecificTags); err != nil {
		return diag.FromErr(err)
	}

	// TODO: In the future, our API or tags need to dynamically grab these values.
	//  Till then, all imports require a destroy + create.
	//_ = d.Set("type", foundrole.RoleType)
	//_ = d.Set("include_default_policies", foundrole.InclDefaultPolicies)

	return nil
}

func resourceAlksIamRoleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS IAM Role Update")

	// enable partial state mode
	d.Partial(true)

	if d.HasChange("enable_alks_access") {
		// try updating enable_alks_access
		if err := updateAlksAccess(d, meta); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags_all") {
		// try updating enable_alks_access
		if err := updateIamTags(d, meta); err != nil {
			return diag.FromErr(err)
		}
	}

	d.Partial(false)

	return resourceAlksIamRoleRead(ctx, d, meta)
}

func updateAlksAccess(d *schema.ResourceData, meta interface{}) error {
	var alksAccess = d.Get("enable_alks_access").(bool)
	var roleArn = d.Get("arn").(string)
	providerStruct := meta.(*AlksClient)
	client := providerStruct.client
	if err := validateIAMEnabled(client); err != nil {
		return err
	}
	// create the machine identity
	if alksAccess {
		_, err := client.AddRoleMachineIdentity(roleArn)
		if err != nil {
			return err
		}
	} else {
		// delete the machine identity
		_, err := client.DeleteRoleMachineIdentity(roleArn)
		if err != nil {
			return err
		}
	}
	return nil
}

func updateIamTags(d *schema.ResourceData, meta interface{}) error {
	providerStruct := meta.(*AlksClient)
	client := providerStruct.client

	if err := validateIAMEnabled(client); err != nil {
		return err
	}

	//Do a read to get existing tags.  If any of those are in ignore_tags, then they are externally managed
	//and they should be included in the update so they don't get removed.
	foundRole, err := client.GetIamRole(d.Id())

	if err != nil {
		return err
	}

	existingTags := tagSliceToMap(foundRole.Tags)
	externalTags := getExternalyManagedTags(existingTags, *providerStruct.ignoreTags)
	internalTags := d.Get("tags_all").(map[string]interface{})

	//Tags includes default tags, role specific tags, and tags that exist externally on the role itself and are specified in ignored_tags
	tags := tagMapToSlice(combineTagMaps(internalTags, externalTags))

	options := alks.UpdateIamRoleRequest{
		RoleName: &foundRole.RoleName,
		Tags:     &tags,
	}

	if _, err := client.UpdateIamRole(&options); err != nil {
		return err
	}
	return nil
}

func migrateState(version int, state *terraform.InstanceState, meta interface{}) (*terraform.InstanceState, error) {
	switch version {
	case 0:
		log.Println("[INFO] Found Instance State v0, migrating to v1")
		return migrateV0toV1(state)
	default:
		return state, fmt.Errorf("Unrecognized version '%d' in schema for instance of ALKS IAM role '%s'", version, state.Attributes["name"])
	}
}

func migrateV0toV1(state *terraform.InstanceState) (*terraform.InstanceState, error) {
	if state.Empty() {
		log.Println("[DEBUG] Empty InstanceState, nothing to migrate")
		return state, nil
	}

	if _, ok := state.Attributes["enable_alks_access"]; !ok {
		log.Printf("[DEBUG] Attributes before migration: %#v", state.Attributes)
		state.Attributes["enable_alks_access"] = "false"
		log.Printf("[DEBUG] Attributes after migration: %#v", state.Attributes)
	}

	return state, nil
}
