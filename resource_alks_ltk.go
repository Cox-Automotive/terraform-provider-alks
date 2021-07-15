package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlksLtk() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAlksLtkCreate,
		ReadContext:   resourceAlksLtkRead,
		DeleteContext: resourceAlksLtkDelete,
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

func resourceAlksLtkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS LTK User Create")

	var iamUsername = d.Get("iam_username").(string)

	client := meta.(*alks.Client)
	if err := validateIAMEnabled(client); err != nil {
		return diag.FromErr(err)
	}

	resp, err := client.CreateLongTermKey(iamUsername)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(iamUsername)
	_ = d.Set("iam_user_arn", resp.IAMUserArn)
	_ = d.Set("access_key", resp.AccessKey)
	_ = d.Set("secret_key", resp.SecretKey)

	log.Printf("[INFO] alks_ltk.id: %v", d.Id())

	return resourceAlksLtkRead(ctx, d, meta)
}

func resourceAlksLtkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS LTK User Read")

	client := meta.(*alks.Client)

	// Check if role exists.
	if d.Id() == "" || d.Id() == "none" {
		return nil
	}

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

func resourceAlksLtkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[INFO] ALKS LTK User Delete")

	client := meta.(*alks.Client)
	if err := validateIAMEnabled(client); err != nil {
		return diag.FromErr(err)
	}

	if _, err := client.DeleteLongTermKey(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
