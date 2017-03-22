package main

import (
	"log"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlksSession() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlksSessionCreate,
		Read:   resourceAlksSessionRead,
		Update: resourceAlksSessionUpdate,
		Delete: resourceAlksSessionDelete,

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

// No-op: ALKS Session is currently read-only
func resourceAlksSessionCreate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlksSessionRead(d, meta)
}

func resourceAlksSessionRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS Session Read")

	var useIam = d.Get("use_iam").(bool)
	var name = d.Get("name").(string)
	var duration int = 2
	if useIam {
		duration = 1
	}

	client := meta.(*alks.Client)
	resp, err := client.CreateSession(duration, useIam)

	if err != nil {
		return err
	}

	d.SetId(name)
	d.Set("access_key", resp.AccessKey)
	d.Set("secret_key", resp.SecretKey)
	d.Set("session_token", resp.SessionToken)

	log.Printf("[INFO] alks_session.id: %v", d.Id())

	return nil
}

// No-op: ALKS Session is currently read-only
func resourceAlksSessionUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlksSessionRead(d, meta)
}

// No-op: ALKS Session is currently read-only
func resourceAlksSessionDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
