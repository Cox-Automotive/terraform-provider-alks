package main

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"alks": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("Error: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ALKS_URL"); v == "" {
		t.Fatal("ALKS_URL must be set for acceptance tests")
	}
	if v := os.Getenv("AWS_ACCESS_KEY_ID"); v == "" {
		t.Fatal("AWS_ACCESS_KEY_ID must be set for acceptance tests")
	}
	if v := os.Getenv("AWS_SECRET_ACCESS_KEY"); v == "" {
		t.Fatal("AWS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("AWS_SESSION_TOKEN"); v == "" {
		t.Fatal("AWS_SESSION_TOKEN must be set for acceptance tests")
	}
}

func testAccProvider_DefaultTags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testProviderWithDefaultTags,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIgnoreTagsKeys(testAccProviders["alks"], []string{"defaultTagKey1", "defaultTagKey2"}),
				),
			},
		},
	})
}

func testAccCheckIgnoreTagsKeys(providers *schema.Provider, expectedKeys []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerStruct := meta.(*AlksClient)
		client := providerStruct.client

		defaultTags := providerStruct.defaultTags
		found := false
		for _, k := range expectedKeys {
			for _, i := range defaultTags {
				if k == i.Key {
					found = true
					break
				}
			}
			if !found {
				return diag.FromErr("Expected Key Not found")
			}
		}
		return nil
	}
}

const testProviderWithDefaultTags = `
	provider "alks" {
		default_tags {
			tags = {
				defaultTagKey1 = "defaultTagValue1"
				defaultTagKey2 = "defaultTagValue2"
			}
		}
	}
`
