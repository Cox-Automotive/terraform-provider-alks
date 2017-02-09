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
	if v := os.Getenv("ALKS_ROLE"); v == "" {
		t.Fatal("ALKS_ROLE must be set for acceptance tests")
	}
	if v := os.Getenv("ALKS_ACCOUNT"); v == "" {
		t.Fatal("ALKS_ACCOUNT must be set for acceptance tests")
	}
	if v := os.Getenv("ALKS_URL"); v == "" {
		t.Fatal("ALKS_URL must be set for acceptance tests")
	}
	if v := os.Getenv("ALKS_USERNAME"); v == "" {
		t.Fatal("ALKS_USERNAME must be set for acceptance tests")
	}
	if v := os.Getenv("ALKS_PASSWORD"); v == "" {
		t.Fatal("ALKS_PASSWORD must be set for acceptance tests")
	}
}
