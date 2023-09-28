package main

import (
	"context"

	"github.com/Cox-Automotive/alks-go"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlksSession() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAlksSessionRead,
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_key": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"session_token": {
				Type:     schema.TypeString,
				Computed: true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceAlksSessionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerStruct := meta.(*AlksClient)
	client := providerStruct.client
	credentials := client.Credentials.(interface{})

	// Type assertion: Since STS is a struct, we need to cast it to the interface.
	stsCredentials, ok := credentials.(*alks.STS)

	if ok {
		_ = d.Set("access_key", stsCredentials.AccessKey)
		_ = d.Set("secret_key", stsCredentials.SecretKey)
		_ = d.Set("session_token", stsCredentials.SessionToken)
		// When the below line was not present I could not get credentials to be returned
		// as terraform outputs during runtime.
		d.SetId(client.AccountDetails.Account)
	}

	// If not ok, we return empty values for credentials and move along.  Is this what we want?

	return nil
}
