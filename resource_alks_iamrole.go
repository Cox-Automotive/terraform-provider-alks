package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	alks "github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func resourceAlksIamRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlksIamRoleCreate,
		Read:   resourceAlksIamRoleRead,
		Update: resourceAlksIamRoleUpdate,
		Exists: resourceAlksIamRoleExists,
		Delete: resourceAlksIamRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		MigrateState:  migrateState,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"include_default_policies": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"role_added_to_ip": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_alks_access": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		},
	}
}

func resourceAlksIamTrustRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlksIamTrustRoleCreate,
		Read:   resourceAlksIamRoleRead,
		Update: resourceAlksIamRoleUpdate,
		Exists: resourceAlksIamRoleExists,
		Delete: resourceAlksIamRoleDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		SchemaVersion: 1,
		MigrateState:  migrateState,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"trust_arn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"role_added_to_ip": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_alks_access": &schema.Schema{
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
			},
		},
	}
}

func resourceAlksIamRoleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS IAM Role Create")

	var roleName = d.Get("name").(string)
	var roleType = d.Get("type").(string)
	var incDefPol = d.Get("include_default_policies").(bool)
	var enableAlksAccess = d.Get("enable_alks_access").(bool)

	client := meta.(*alks.Client)
	resp, err := client.CreateIamRole(roleName, roleType, incDefPol, enableAlksAccess)

	if err != nil {
		return err
	}

	d.SetId(resp.RoleName)

	log.Printf("[INFO] alks_iamrole.id: %v", d.Id())

	return resourceAlksIamRoleRead(d, meta)
}

func resourceAlksIamTrustRoleCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS IAM Trust Role Create")

	var roleName = d.Get("name").(string)
	var roleType = d.Get("type").(string)
	var trustArn = d.Get("trust_arn").(string)
	var enableAlksAccess = d.Get("enable_alks_access").(bool)

	client := meta.(*alks.Client)

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
		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}

	response := *resp

	d.SetId(response.RoleName)

	log.Printf("[INFO] alks_iamtrustrole.id: %v", d.Id())

	return resourceAlksIamRoleRead(d, meta)
}

func resourceAlksIamRoleDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS IAM Role Delete")

	client := meta.(*alks.Client)
	err := client.DeleteIamRole(d.Id())

	if err != nil {
		return err
	}

	return nil
}

func resourceAlksIamRoleExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	log.Printf("[INFO] ALKS IAM Role Exists")

	client := meta.(*alks.Client)

	foundRole, err := client.GetIamRole(d.Id())

	if err != nil {
		// TODO: Clean-up this logic, likely by improving the error responses from `alks-go`
		if strings.Contains(err.Error(), "Role not found") {
			return false, nil
		}

		return false, err
	}

	if foundRole == nil {
		return false, nil
	}

	return true, nil
}

func resourceAlksIamRoleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS IAM Role Read")

	client := meta.(*alks.Client)
	foundRole, err := client.GetIamRole(d.Id())

	if err != nil {
		d.SetId("")
		return err
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

func resourceAlksIamRoleUpdate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS IAM Role Update")

	// enable partial state mode
	d.Partial(true)

	if d.HasChange("enable_alks_access") {
		// try updating enable_alks_access
		if err := updateAlksAccess(d, meta); err != nil {
			return err
		}

		d.SetPartial("enable_alks_access")
	}

	d.Partial(false)

	return nil
}

func updateAlksAccess(d *schema.ResourceData, meta interface{}) error {
	var alksAccess = d.Get("enable_alks_access").(bool)
	var roleArn = d.Get("arn").(string)
	client := meta.(*alks.Client)
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
