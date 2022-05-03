package main

import (
	"context"
	"log"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mitchellh/go-homedir"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "This is the base URL to ALKS service. It must be provided, but it can also be sourced from the ALKS_URL environment variable.",
				DefaultFunc: schema.EnvDefaultFunc("ALKS_URL", nil),
			},
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the AWS access key. It must be provided, but it can also be sourced from the ALKS_ACCESS_KEY_ID or AWS_ACCESS_KEY_ID environment variable.",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_ACCESS_KEY_ID",
					"ALKS_ACCESS_KEY_ID",
				}, nil),
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the AWS secret key. It must be provided, but it can also be sourced from the ALKS_SECRET_ACCESS_KEY or AWS_SECRET_ACCESS_KEY environment variable",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_SECRET_ACCESS_KEY",
					"ALKS_SECRET_ACCESS_KEY",
				}, nil),
			},
			"token": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "This is the AWS session token. It must be provided, but it can also be sourced from the ALKS_SESSION_TOKEN or AWS_SESSION_TOKEN environment variable",
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_SESSION_TOKEN",
					"ALKS_SESSION_TOKEN",
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
			"account": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The account which you'd like to retrieve credentials for.",
				DefaultFunc: schema.EnvDefaultFunc("Account", nil),
			},
			"role": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The role which you'd like to retrieve credentials for.",
				DefaultFunc: schema.EnvDefaultFunc("Role", nil),
			},
			"assume_role":  assumeRoleSchema(),
			"default_tags": defaultTagsSchema(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"alks_iamrole": resourceAlksIamRole(),
			// "alks_iamtrustrole": resourceAlksIamTrustRole(),
			"alks_ltk": resourceAlksLtk(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"alks_keys": dataSourceAlksKeys(),
		},

		ConfigureContextFunc: providerConfigure,
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

func defaultTagsSchema() *schema.Schema {
	return &schema.Schema{
		Type:        schema.TypeList,
		Optional:    true,
		MaxItems:    1,
		Description: "Configuration block with settings to default resource tags across all resources.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"tags": {
					Type:        schema.TypeMap,
					Optional:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Description: "Resource tags to default across all resources",
				},
			},
		},
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {

	var diags diag.Diagnostics

	config := Config{
		URL:       d.Get("url").(string),
		AccessKey: d.Get("access_key").(string),
		SecretKey: d.Get("secret_key").(string),
		Token:     d.Get("token").(string),
		Profile:   d.Get("profile").(string),
		Account:   d.Get("account").(string),
		Role:      d.Get("role").(string),
		// DefaultTags: expandProviderDefaultTags(d.Get("default_tags").([]interface{})),
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
		return nil, diag.FromErr(err)
	}
	config.CredsFilename = credsPath
	defaultTags := expandProviderDefaultTags(d.Get("default_tags").([]interface{}))

	c, err := config.Client()
	if err != nil {
		return nil, diag.FromErr(err)
	}

	alksClient := &AlksClient{}
	alksClient.client = c
	if defaultTags != nil {
		alksClient.defaultTags = defaultTags
	}

	log.Println("[INFO] Initializing ALKS client")
	return alksClient, diags
	// return c, diags
}

func expandProviderDefaultTags(l []interface{}) []alks.Tag {
	if len(l) == 0 || l[0] == nil {
		return nil
	}

	m := l[0].(map[string]interface{})
	tagSlice := tagMapToSlice(m["tags"].(map[string]interface{}))

	return tagSlice
}

type AlksClient struct {
	client      *alks.Client
	defaultTags []alks.Tag //Not making this a pointer because I was having to check everywhere if it was nil
}
