package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
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
			},
			"name_prefix": {
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ForceNew:      true,
				ConflictsWith: []string{"name"},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"include_default_policies": {
				Type:     schema.TypeBool,
				Required: true,
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
			},
			"tags":     TagsSchema(),
			"tags_all": TagsSchemaComputed(),
		},
		CustomizeDiff: SetTagsDiff,
	}
}

func TagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func TagsSchemaComputed() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeMap,
		Optional: true,
		Computed: true,
		Elem:     &schema.Schema{Type: schema.TypeString},
	}
}

func resourceAlksIamRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS IAM Role Create")
	fmt.Printf("In Role Create\n")
	var roleName = NameWithPrefix(d.Get("name").(string), d.Get("name_prefix").(string))
	var roleType = d.Get("type").(string)
	var incDefPol = d.Get("include_default_policies").(bool)
	var enableAlksAccess = d.Get("enable_alks_access").(bool)
	var rawTemplateFields = d.Get("template_fields").(map[string]interface{})
	var maxSessionDurationInSeconds = d.Get("max_session_duration_in_seconds").(int)
	var tags = tagMapToSlice(d.Get("tags").(map[string]interface{}))

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

	defaultTags := []alks.Tag{}
	if (*providerStruct).defaultTags != nil {
		defaultTags = (*providerStruct).defaultTags
	}
	allTags := combineTagsWithDefault(tags, defaultTags)
	//  client := meta.(*alks.Client)

	if err := validateIAMEnabled(client); err != nil {
		return diag.FromErr(err)
	}

	options := &alks.CreateIamRoleOptions{
		RoleName:                    &roleName,
		RoleType:                    &roleType,
		IncludeDefaultPolicies:      &include,
		AlksAccess:                  &enableAlksAccess,
		TemplateFields:              &templateFields,
		MaxSessionDurationInSeconds: &maxSessionDurationInSeconds,
		Tags:                        &allTags,
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
	fmt.Printf("In Role Read\n")
	providerStruct := meta.(*AlksClient)
	client := providerStruct.client

	defaultTags := providerStruct.defaultTags

	// Check if role exists.
	if d.Id() == "" || d.Id() == "none" {
		return nil
	}

	foundRole, err := client.GetIamRole(d.Id())

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] alks_iamrole.id %v", d.Id())

	_ = d.Set("name", foundRole.RoleName)
	_ = d.Set("name_prefix", NamePrefixFromName(foundRole.RoleName))
	_ = d.Set("arn", foundRole.RoleArn)
	_ = d.Set("ip_arn", foundRole.RoleIPArn)
	_ = d.Set("enable_alks_access", foundRole.AlksAccess)

	allTags := foundRole.Tags
	for _, t := range allTags {
		fmt.Printf("Read Tag: %s\n", t.Key)
	}
	roleSpecificTags := removeDefaultTags(tagSliceToMap(allTags), defaultTags)

	if err := d.Set("tags", tagSliceToMap(roleSpecificTags)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags_all", tagSliceToMap(allTags)); err != nil {
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
	fmt.Printf("In Role Update\n")

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
	fmt.Printf("In updateIamTags\n")
	providerStruct := meta.(*AlksClient)
	client := providerStruct.client

	if err := validateIAMEnabled(client); err != nil {
		return err
	}

	tags := tagMapToSlice(d.Get("tags_all").(map[string]interface{}))
	roleName := NameWithPrefix(d.Get("name").(string), d.Get("name_prefix").(string))
	options := alks.UpdateIamRoleRequest{
		RoleName: &roleName,
		Tags:     &tags,
	}

	if _, err := client.UpdateIamRole(&options); err != nil {
		return err
	}
	return nil
}

//Combines tags defined on an individual resource with the default tags listed on the provider block
//Resource specific tags will overwrite default tags
func combineTagsWithDefault(tags []alks.Tag, defaultTags []alks.Tag) []alks.Tag {
	defaultTagsMap := tagSliceToMap(defaultTags)

	for _, t := range tags {
		defaultTagsMap[t.Key] = t.Value
	}
	allTags := tagMapToSlice(defaultTagsMap)

	return allTags
}

//Removes default tags from a map of role specific + default tags
func removeDefaultTags(allTags map[string]interface{}, defalutTags []alks.Tag) []alks.Tag {
	for _, t := range defalutTags {
		delete(allTags, t.Key)
	}

	return tagMapToSlice(allTags)
}

func tagMapToSlice(tagMap map[string]interface{}) []alks.Tag {
	tags := []alks.Tag{}
	for k, v := range tagMap {
		tag := alks.Tag{Key: k, Value: v.(string)}
		tags = append(tags, tag)
	}
	return tags
}

func tagSliceToMap(tagSlice []alks.Tag) map[string]interface{} {
	tagMap := make(map[string]interface{})
	for _, t := range tagSlice {
		tagMap[t.Key] = t.Value
	}
	return tagMap
}

func SetTagsDiff(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	defaultTags := meta.(*AlksClient).defaultTags

	resourceTags := tagMapToSlice(diff.Get("tags").(map[string]interface{}))

	allTags := combineTagsWithDefault(resourceTags, defaultTags)

	// To ensure "tags_all" is correctly computed, we explicitly set the attribute diff
	// when the merger of resource-level tags onto provider-level tags results in n > 0 tags,
	// otherwise we mark the attribute as "Computed" only when their is a known diff (excluding an empty map)
	// or a change for "tags_all".

	if len(allTags) > 0 {
		if err := diff.SetNew("tags_all", tagSliceToMap(allTags)); err != nil {
			return fmt.Errorf("error setting new tags_all diff: %w", err)
		}
	} else if len(diff.Get("tags_all").(map[string]interface{})) > 0 {
		if err := diff.SetNewComputed("tags_all"); err != nil {
			return fmt.Errorf("error setting tags_all to computed: %w", err)
		}
	} else if diff.HasChange("tags_all") {
		if err := diff.SetNewComputed("tags_all"); err != nil {
			return fmt.Errorf("error setting tags_all to computed: %w", err)
		}
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
