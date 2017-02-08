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
                Description: "Base URL to ALKS service",
                DefaultFunc: schema.EnvDefaultFunc("ALKS_URL", nil),
            },
            "username": &schema.Schema{
                Type:        schema.TypeString,
                Required:    true,
                Description: "Username used to login to ALKS",
                DefaultFunc: schema.EnvDefaultFunc("ALKS_USERNAME", nil),
            },
            "password": &schema.Schema{
                Type:        schema.TypeString,
                Required:    true,
                Description: "Password used to login to ALKS",
                DefaultFunc: schema.EnvDefaultFunc("ALKS_PASSWORD", nil),
            },
            "account": &schema.Schema{
                Type:        schema.TypeString,
                Required:    true,
                Description: "ALKS Account to use",
                DefaultFunc: schema.EnvDefaultFunc("ALKS_ACCOUNT", nil),
            },
            "role": &schema.Schema{
                Type:        schema.TypeString,
                Required:    true,
                Description: "ALKS role to use",
                DefaultFunc: schema.EnvDefaultFunc("ALKS_ROLE", nil),
            },
        },

        ResourcesMap: map[string]*schema.Resource{
            "alks_iamrole": resourceAlksIamRole(),
        },

        ConfigureFunc: providerConfigure,
    }
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
    config := Config{
        Url:      d.Get("url").(string),
        Username: d.Get("username").(string),
        Password: d.Get("password").(string),
        Account:  d.Get("account").(string),
        Role:     d.Get("role").(string),
    }

    log.Println("[INFO] Initializing ALKS client")
    return config.Client()
}
