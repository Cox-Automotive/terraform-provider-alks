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
				Type:     schema.TypeString,
				Computed: true,
			},
			"session_token": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlksSessionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	providerStruct := meta.(*AlksClient)
	client := providerStruct.client
	credentials := client.Credentials.(interface{})

	stsCredentials, ok := credentials.(alks.STS)

	if ok {
		_ = d.Set("access_key", stsCredentials.AccessKey)
		_ = d.Set("secret_key", stsCredentials.SecretKey)
		_ = d.Set("session_token", stsCredentials.SessionToken)
	}

	return nil
}
