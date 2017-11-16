package main

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"alks": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("Error: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("ALKS_URL"); v == "" {
		t.Fatal("ALKS_URL must be set for acceptance tests")
	}
	if v := os.Getenv("ALKS_ACCESS_KEY_ID"); v == "" {
		t.Fatal("ALKS_ACCESS_KEY_ID must be set for acceptance tests")
	}
	if v := os.Getenv("ALKS_SECRET_ACCESS_KEY"); v == "" {
		t.Fatal("ALKS_SECRET_ACCESS_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("ALKS_SESSION_TOKEN"); v == "" {
		t.Fatal("ALKS_SESSION_TOKEN must be set for acceptance tests")
	}
}
