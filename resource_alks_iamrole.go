package main

import (
	"context"
	"fmt"
	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"log"
	"strings"
	"time"
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
		MigrateState: migrateState,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
		},
	}
}

func resourceAlksIamTrustRole() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlksIamTrustRoleCreate,
		ReadContext:   resourceAlksIamRoleRead,
		UpdateContext: resourceAlksIamRoleUpdate,
		DeleteContext: resourceAlksIamRoleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		SchemaVersion: 1,
		MigrateState: migrateState,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"trust_arn": {
				Type:     schema.TypeString,
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
		},
	}
}

func resourceAlksIamRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS IAM Role Create")

	var roleName = d.Get("name").(string)
	var roleType = d.Get("type").(string)
	var incDefPol = d.Get("include_default_policies").(bool)
	var enableAlksAccess = d.Get("enable_alks_access").(bool)
	var rawTemplateFields = d.Get("template_fields").(map[string]interface{})

	templateFields := make(map[string]string)
	for k, v := range rawTemplateFields {
		templateFields[k] = v.(string)
	}

	client := meta.(*alks.Client)
	if err := validateIAMEnabled(client); err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.CreateIamRole(roleName, roleType, templateFields, incDefPol, enableAlksAccess)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resp.RoleName)
	_ = d.Set("role_added_to_ip", resp.RoleAddedToIP)

	log.Printf("[INFO] alks_iamrole.id: %v", d.Id())

	return resourceAlksIamRoleRead(ctx, d, meta)
}

func resourceAlksIamTrustRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS IAM Trust Role Create")

	var roleName = d.Get("name").(string)
	var roleType = d.Get("type").(string)
	var trustArn = d.Get("trust_arn").(string)
	var enableAlksAccess = d.Get("enable_alks_access").(bool)

	client := meta.(*alks.Client)
	if err := validateIAMEnabled(client); err != nil {
		return diag.FromErr(err)
	}

	var resp *alks.IamRoleResponse
	err := resource.Retry(2*time.Minute, func() *resource.RetryError {
		var err error
		resp, err = client.CreateIamTrustRole(roleName, roleType, trustArn, enableAlksAccess)
		if err != nil {
			if strings.Contains(err.Error(), "Role already exists") || strings.Contains(err.Error(), "Instance profile exists") {
				return resource.NonRetryableError(err)
			}

			// Due to IAM eventual consistency, you might've just created a role that you need to link to a trust
			// We'll keep checking every 15 seconds for up to 2 minutes to see if the role appears
			time.Sleep(15 * time.Second)
			return resource.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return diag.FromErr(err)
	}

	response := *resp

	d.SetId(response.RoleName)
	_ = d.Set("role_added_to_ip", resp.RoleAddedToIP)

	log.Printf("[INFO] alks_iamtrustrole.id: %v", d.Id())

	return resourceAlksIamRoleRead(ctx, d, meta)
}

func resourceAlksIamRoleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS IAM Role Delete")

	client := meta.(*alks.Client)
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

	client := meta.(*alks.Client)

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
	_ = d.Set("arn", foundRole.RoleArn)
	_ = d.Set("ip_arn", foundRole.RoleIPArn)
	_ = d.Set("enable_alks_access", foundRole.AlksAccess)

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

	d.Partial(false)

	return resourceAlksIamRoleRead(ctx, d, meta)
}

func updateAlksAccess(d *schema.ResourceData, meta interface{}) error {
	var alksAccess = d.Get("enable_alks_access").(bool)
	var roleArn = d.Get("arn").(string)
	client := meta.(*alks.Client)
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