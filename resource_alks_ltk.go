package main

import (
	"log"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlksLtk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlksLtkCreate,
		Read:   resourceAlksLtkRead,
		Delete: resourceAlksLtkDelete,
		Exists: resourceAlksLtkExists,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		SchemaVersion: 1,

		Schema: map[string]*schema.Schema{
			"iam_username": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"iam_user_arn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_key": {
				Sensitive: true,
				Type:      schema.TypeString,
				Computed:  true,
			},
			"secret_key": {
				Sensitive: true,
				Type:      schema.TypeString,
				Computed:  true,
			},
		},
	}
}

func resourceAlksLtkCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS LTK User Create")

	var iamUsername = d.Get("iam_username").(string)

	client := meta.(*alks.Client)
	if err := validateIAMEnabled(client); err != nil {
		return err
	}

	resp, err := client.CreateLongTermKey(iamUsername)
	if err != nil {
		return err
	}

	d.SetId(iamUsername)
	_ = d.Set("iam_user_arn", resp.IAMUserArn)
	_ = d.Set("access_key", resp.AccessKey)
	_ = d.Set("secret_key", resp.SecretKey)

	log.Printf("[INFO] alks_ltk.id: %v", d.Id())

	return resourceAlksLtkRead(d, meta)
}

func resourceAlksLtkRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS LTK User Read")

	client := meta.(*alks.Client)
	resp, err := client.GetLongTermKey(d.Id())

	if err != nil {
		d.SetId("")
		return nil
	}

	log.Printf("[INFO] alks_ltk.id: %v", d.Id())

	_ = d.Set("iam_username", resp.UserName)
	_ = d.Set("access_key", resp.AccessKeyID)

	return nil
}

func resourceAlksLtkDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS LTK User Delete")

	client := meta.(*alks.Client)
	if err := validateIAMEnabled(client); err != nil {
		return err
	}

	if _, err := client.DeleteLongTermKey(d.Id()); err != nil {
		return err
	}

	return nil
}

func resourceAlksLtkExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Printf("[INFO] ALKS LTK User Exists")

	client := meta.(*alks.Client)
	resp, err := client.GetLongTermKey(d.Id())

	if err != nil {
		return false, err
	}

	// We can get a 200, but an empty string so this is the condition to check for.
	if len(resp.LongTermKey.UserName) == 0 {
		return false, nil
	}

	return true, nil
}
