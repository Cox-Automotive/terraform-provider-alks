package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	homedir "github.com/mitchellh/go-homedir"
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
				Optional:    true,
				Description: "This is the AWS access key. It must be provided, but it can also be sourced from the ALKS_ACCESS_KEY_ID or AWS_ACCESS_KEY_ID environment variable.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"ALKS_ACCESS_KEY_ID",
					"AWS_ACCESS_KEY_ID",
				}, nil),
			},
			"secret_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the AWS secret key. It must be provided, but it can also be sourced from the ALKS_SECRET_ACCESS_KEY or AWS_SECRET_ACCESS_KEY environment variable",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"ALKS_SECRET_ACCESS_KEY",
					"AWS_SECRET_ACCESS_KEY",
				}, nil),
			},
			"token": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the AWS session token. It must be provided, but it can also be sourced from the ALKS_SESSION_TOKEN or AWS_SESSION_TOKEN environment variable",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"ALKS_SESSION_TOKEN",
					"AWS_SESSION_TOKEN",
				}, nil),
			},
			"profile": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The profile for API operations. Used in place of STS tokens.",
				DefaultFunc: schema.EnvDefaultFunc("AWS_PROFILE", nil),
			},
			"shared_credentials_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The path to the shared credentials file. If not set this defaults to ~/.aws/credentials.",
				DefaultFunc: schema.EnvDefaultFunc("AWS_SHARED_CREDENTIALS_FILE", nil),
			},
			"assume_role": assumeRoleSchema(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"alks_iamrole":          resourceAlksIamRole(),
			"alks_iamtrustrole":     resourceAlksIamTrustRole(),
			"alks_machine_identity": resourceAlksIamMachineIdentity(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func assumeRoleSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"role_arn": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "(Required) Role ARN to assume before calling ALKS",
				},
				"session_name": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "(Optional) Session name to use when making the AssumeRole call.  See AWS SDK for more details.",
				},
				"external_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "(Optional) The external ID to use when making the AssumeRole call. See AWS SDK for more details.",
				},
				"policy": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "(Optional) Additional policy restrictions to apply to the result STS session.  See AWS SDK for more details.",
				},
			},
		},
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		URL:       d.Get("url").(string),
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
		Token:     d.Get("token").(string),
		Profile:   d.Get("profile").(string),
	}

	assumeRoleList := d.Get("assume_role").(*schema.Set).List()
	if len(assumeRoleList) == 1 {
		assumeRole := assumeRoleList[0].(map[string]interface{})
		config.AssumeRole.RoleARN = assumeRole["role_arn"].(string)
		config.AssumeRole.SessionName = assumeRole["session_name"].(string)
		config.AssumeRole.ExternalID = assumeRole["external_id"].(string)
		config.AssumeRole.Policy = assumeRole["policy"].(string)
	}

	// Set CredsFilename, expanding home directory
	credsPath, err := homedir.Expand(d.Get("shared_credentials_file").(string))
	if err != nil {
		return nil, err
	}
	config.CredsFilename = credsPath

	log.Println("[INFO] Initializing ALKS client")
	return config.Client()
}
