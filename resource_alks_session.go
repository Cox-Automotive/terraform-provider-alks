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
		Exists: nil,
		Delete: resourceAlksSessionDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"use_iam": &schema.Schema{
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
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

func resourceAlksSessionRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS Session Read")

	return nil
}

func resourceAlksSessionDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS Session Delete")

	return nil
}
