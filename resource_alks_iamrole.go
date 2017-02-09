package main

import (
    "log"

    "github.com/hashicorp/terraform/helper/schema"
)

func resourceAlksIamRole() *schema.Resource {
    return &schema.Resource{
        Create: resourceAlksIamRoleCreate,
        Read:   resourceAlksIamRoleRead,
        Delete: resourceAlksIamRoleDelete,

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
                Type:     schema.TypeString,
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
        },
    }
}

func resourceAlksIamRoleCreate(d *schema.ResourceData, meta interface{}) error {
    log.Printf("[DEBUG] ALKS IAM Role Create")

    var roleName  = d.Get("name").(string)
    var roleType  = d.Get("type").(string)
    var incDefPol = d.Get("include_default_policies").(bool)

    client := meta.(*AlksClient)
    resp, err := client.CreateIamRole(roleName, roleType, incDefPol)

    if err != nil {
        return err
    }

    d.SetId(resp.RoleArn)
    d.Set("arn", resp.RoleArn)
    d.Set("ip_arn", resp.RoleIPArn)
    d.Set("role_added_to_ip", resp.RoleAddedToIP)

    log.Printf("[INFO] Created role: %s with ARN: %s", resp.RoleName, resp.RoleArn)

    return nil
}

func resourceAlksIamRoleDelete(d *schema.ResourceData, meta interface{}) error {
    log.Printf("[DEBUG] ALKS IAM Role Delete")

    return nil
}

func resourceAlksIamRoleRead(d *schema.ResourceData, meta interface{}) error {
    log.Printf("[DEBUG] ALKS IAM Role Read")

    return nil
}