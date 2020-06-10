package main

import (
	alks "github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
)

// TODO: Find a way to NOT store 'secret_key' to TF state file.
func resourceAlksLtk() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlksLtkCreate,
		Read:   resourceAlksLtkRead,
		Delete: resourceAlksLtkDelete,

		SchemaVersion: 1,
		MigrateState:  migrateState,

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
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"ltks": &schema.Schema{
				Type:     schema.TypeMap,
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

	d.Set("account_id", resp.AccountDetails.Account)
	d.Set("role_name", resp.AccountDetails.Role)

	d.Set("role", resp.Role)
	d.Set("action", resp.Action)
	d.Set("added_iam_user_to_group", resp.AddedIAMUserToGroup)
	d.Set("partial_error", resp.PartialError)
	d.Set("iam_user_arn", resp.IAMUserArn)
	d.Set("access_key", resp.AccessKey)
	d.Set("secret_key", resp.SecretKey)

	log.Printf("[INFO] alks_ltk.id: %v", d.Id())

	return nil
}

func resourceAlksLtkRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS LTK Users Read")

	var accountID = d.Get("account_id").(string)
	var roleName = d.Get("role_name").(string)

	client := meta.(*alks.Client)
	resp, err := client.GetLongTermKeys(accountID, roleName)

	if err != nil {
		return err
	}

	d.SetId(resp.RequestID)
	d.Set("ltks", resp.LongTermKeys)

	log.Printf("[INFO] alks_ltk.id: %v", d.Id())

	return nil
}

func resourceAlksLtkDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS LTK User Delete")

	var iamUsername = d.Get("iam_username").(string)

	client := meta.(*alks.Client)
	resp, err := client.DeleteLongTermKey(iamUsername)

	if err != nil {
		return err
	}

	d.Set("role", resp.Role)
	d.Set("action", resp.Action)
	d.Set("added_iam_user_to_group", resp.AddedIAMUserToGroup)
	d.Set("partial_error", resp.PartialError)

	log.Printf("[INFO] alks_ltk.id: %v", d.Id())

	return nil
}
