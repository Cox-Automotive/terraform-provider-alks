package main

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlksKeys() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlksKeysRead,
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"session_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlksKeysRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] ALKS Keys Data Source Read")

	providerStruct := meta.(*AlksClient)
	client := providerStruct.client
	resp, err := client.CreateIamSession()

	if err != nil {
		return err
	}

	// Return the information to user.
	_ = d.Set("access_key", resp.AccessKey)
	_ = d.Set("secret_key", resp.SecretKey)
	_ = d.Set("session_token", resp.SessionToken)
	_ = d.Set("account", client.AccountDetails.Account)
	_ = d.Set("role", strings.Split(client.AccountDetails.Role, "/")[0])

	d.SetId(client.AccountDetails.Account)

	return nil
}
