package main

import (
    "log"

    "github.com/Cox-Automotive/alks-go"
    "github.com/hashicorp/terraform/helper/schema"
)

func resourceAlksSession() *schema.Resource {
    return &schema.Resource{
        Create: resourceAlksSessionCreate,
        Read:   nil,
        Exists: nil,
        Delete: nil,

        Schema: map[string]*schema.Schema{
            "use_iam": &schema.Schema{
                Type:     schema.TypeBool,
                Required: true,
            },
            "access_key": &schema.Schema{
                Type:     schema.TypeString,
                Computed: true,
            },
            "secret_key": &schema.Schema{
                Type:     schema.TypeString,
                Computed: true,
            },
            "session_token": &schema.Schema{
                Type:     schema.TypeString,
                Computed: true,
            },
        },
    }
}

func resourceAlksSessionCreate(d *schema.ResourceData, meta interface{}) error {
    log.Printf("[INFO] ALKS Session Create")

    return nil
}