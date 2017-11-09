package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the base URL to ALKS service. It must be provided, but it can also be sourced from the ALKS_URL environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("ALKS_URL", nil),
			},
			"access_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the AWS access key. It must be provided, but it can also be sourced from the ALKS_ACCESS_KEY_ID environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("ALKS_ACCESS_KEY_ID", nil),
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the AWS secret key. It must be provided, but it can also be sourced from the ALKS_SECRET_ACCESS_KEY environment variable",
				DefaultFunc: schema.EnvDefaultFunc("ALKS_SECRET_ACCESS_KEY", nil),
			},
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the AWS session token. It must be provided, but it can also be sourced from the ALKS_SESSION_TOKEN environment variable",
				DefaultFunc: schema.EnvDefaultFunc("ALKS_SESSION_TOKEN", nil),
			},
			"account": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the ALKS Account to use. It must be provided, but it can also be sourced from the ALKS_ACCOUNT environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("ALKS_ACCOUNT", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"alks_iamrole": resourceAlksIamRole(),
			"alks_session": resourceAlksSession(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Url:       d.Get("url").(string),
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
		Token:     d.Get("token").(string),
		Account:   d.Get("account").(string),
	}

	log.Println("[INFO] Initializing ALKS client")
	return config.Client()
}
