package main

import (
	alks "github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

func resourceAlksLtk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlksLtkCreate,
		Read:   resourceAlksLtkRead,
		Delete: resourceAlksLtkDelete,

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"iam_username": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_name": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"account": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"action": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"added_iam_user_to_group": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"partial_error": &schema.Schema{
				Type:     schema.TypeBool,
				Computed: true,
			},
			"iam_user_arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_key": &schema.Schema{
				Sensitive: true,
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_key": &schema.Schema{
				Sensitive: true,
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlksLtkCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS LTK User Create")

	var iamUsername = d.Get("iam_username").(string)

	client := meta.(*alks.Client)
	resp, err := client.CreateLongTermKey(iamUsername)

	if err != nil {
		return err
	}

	d.SetId(iamUsername)
	_ = d.Set("iam_user_arn", resp.IAMUserArn)
	_ = d.Set("access_key", resp.AccessKey)
	_ = d.Set("secret_key", resp.SecretKey)

	log.Printf("[INFO] alks_ltk.id: %v", d.Id())

	return nil
}

func resourceAlksLtkRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS LTK Users Read")

	client := meta.(*alks.Client)
	resp, err := client.GetLongTermKey(d.Id())

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] alks_ltk.id: %v", d.Id())

	return populateResourceDataFromLTK(resp, d)
}

func resourceAlksLtkDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS LTK User Delete")

	client := meta.(*alks.Client)
	_, err := client.DeleteLongTermKey(d.Id())

	if err != nil {
		return err
	}

	return nil
}

func populateResourceDataFromLTK(longTermKey *alks.GetLongTermKeyResponse, d *schema.ResourceData) error {
	d.SetId(longTermKey.UserName)
	_ = d.Set("access_key", longTermKey.AccessKeyID)
	_ = d.Set("create_date", longTermKey.CreateDate)

	return nil
}
