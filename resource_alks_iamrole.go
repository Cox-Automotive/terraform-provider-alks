package main

import (
    // "fmt"
    // "time"
    "log"
    // "github.com/hashicorp/terraform/helper/resource"
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
        },
    }
}

func resourceAlksIamRoleCreate(d *schema.ResourceData, meta interface{}) error {
    log.Printf("[DEBUG] ALKS IAM Role Create")

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