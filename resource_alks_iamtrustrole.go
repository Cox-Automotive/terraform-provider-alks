package main

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

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

func resourceAlksIamTrustRoleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS IAM Trust Role Create")

	var roleName = NameWithPrefix(d.Get("name").(string), d.Get("name_prefix").(string))
	var roleType = d.Get("type").(string)
	var trustArn = d.Get("trust_arn").(string)
	var enableAlksAccess = d.Get("enable_alks_access").(bool)

	client := meta.(*alks.Client)
	if err := validateIAMEnabled(client); err != nil {
		return diag.FromErr(err)
	}

	var resp *alks.IamRoleResponse
	err := resource.RetryContext(ctx, 2*time.Minute, func() *resource.RetryError {
		var err error
		resp, err = client.CreateIamTrustRole(roleName, roleType, trustArn, enableAlksAccess)
		if err != nil {
			if strings.Contains(err.Error(), "Role already exists") || strings.Contains(err.Error(), "Instance profile exists") {
				return resource.NonRetryableError(err)
			}

			// Amazon IAM utilizes an eventual consistency model:
			//   https://docs.aws.amazon.com/IAM/latest/UserGuide/troubleshoot_general.html#troubleshoot_general_eventual-consistency
			//
			// The newly created IAM role may not exist immediately and could result in dependent
			// resources failing non-deterministically.  Loop for 15 second increments up to 2
			// minutes checking to ensure the resouce was successfully created and is visible.

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
