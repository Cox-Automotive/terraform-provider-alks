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
		Exists: resourceAlksIamRoleExists,
		Delete: resourceAlksIamRoleDelete,

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
				ForceNew: true,
			},
		},
	}
}

func resourceAlksIamTrustRole() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlksIamTrustRoleCreate,
		Read:   resourceAlksIamRoleRead,
		Exists: resourceAlksIamRoleExists,
		Delete: resourceAlksIamRoleDelete,

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
				ForceNew: true,
			},
		},
	}
}

func resourceAlksIamMachineIdentity() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlksIamMachineIdentityCreate,
		Read:   resourceAlksIamMachineIdentityRead,
		Exists: resourceAlksIamMachineIdentityExists,
		Delete: resourceAlksIamMachineIdentityDelete,

		SchemaVersion: 1,
		MigrateState:  migrateState,

		Schema: map[string]*schema.Schema{
			"role_arn": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"machine_identity_arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
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
	d.Set("arn", resp.RoleArn)
	d.Set("ip_arn", resp.RoleIPArn)
	d.Set("role_added_to_ip", resp.RoleAddedToIP)

	log.Printf("[INFO] alks_iamrole.id: %v", d.Id())

	return nil
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

			return resource.RetryableError(err)
		}
		return resource.NonRetryableError(err)
	})

	if err != nil {
		return err
	}

	response := *resp

	d.SetId(response.RoleName)
	d.Set("arn", response.RoleArn)
	d.Set("ip_arn", response.RoleIPArn)
	d.Set("role_added_to_ip", response.RoleAddedToIP)

	log.Printf("[INFO] alks_iamtrustrole.id: %v", d.Id())

	return nil
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

	foundrole, err := client.GetIamRole(d.Id())

	if err != nil {
		// TODO: Clean-up this logic, likely by improving the error responses from `alks-go`
		if strings.Contains(err.Error(), "Role not found") {
			return false, nil
		}

		return false, err
	}

	if foundrole == nil {
		return false, nil
	}

	return true, nil
}

func resourceAlksIamRoleRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS IAM Role Read")

	client := meta.(*alks.Client)

	foundrole, err := client.GetIamRole(d.Id())

	if err != nil {
		return err
	}

	return populateResourceDataFromRole(foundrole, d)
}

func populateResourceDataFromRole(role *alks.IamRoleResponse, d *schema.ResourceData) error {
	d.SetId(role.RoleName)
	d.Set("arn", role.RoleArn)
	d.Set("ip_arn", role.RoleIPArn)
	// role type isnt returned by alks api so this will always false report on a remote state change
	// for more info see issue #125 on ALKS repo
	// d.Set("type", role.RoleType)

	return nil
}

func populateResourceDataFromMI(mi *alks.MachineIdentityResponse, d *schema.ResourceData) error {
	d.SetId(mi.MachineIdentityArn)
	d.Set("machine_identity_arn", mi.MachineIdentityArn)

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

func resourceAlksIamMachineIdentityCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS IAM Machine Identity Create")

	var roleArn = d.Get("role_arn").(string)

	client := meta.(*alks.Client)
	resp, err := client.AddRoleMachineIdentity(roleArn)

	if err != nil {
		return err
	}

	d.SetId(roleArn)
	d.Set("machine_identity_arn", resp.MachineIdentityArn)

	log.Printf("[INFO] alks_machine_identity_arn: %v", d.Get("machine_identity_arn").(string))

	return nil
}

func resourceAlksIamMachineIdentityRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS IAM Machine Identity Read")

	client := meta.(*alks.Client)

	foundMI, err := client.SearchRoleMachineIdentity(d.Id())

	if err != nil {
		return err
	}

	return populateResourceDataFromMI(foundMI, d)
}

func resourceAlksIamMachineIdentityExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	log.Printf("[INFO] ALKS IAM Machine Identity Exists")

	client := meta.(*alks.Client)

	foundMI, err := client.SearchRoleMachineIdentity(d.Id())

	if err != nil {
		if strings.Contains(err.Error(), "Could not find a matching record with the given parameters") {
			return false, nil
		}

		return false, err
	}

	if foundMI == nil {
		return false, nil
	}

	return true, nil
}

func resourceAlksIamMachineIdentityDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS IAM Machine Identity Delete")

	var roleArn = d.Get("role_arn").(string)

	client := meta.(*alks.Client)
	_, err := client.DeleteRoleMachineIdentity(roleArn)

	if err != nil {
		return err
	}

	return nil
}
