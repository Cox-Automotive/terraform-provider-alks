package main

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

func TestExpandProviderDefaultTags_Empty(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"tags": map[string]interface{}{},
		},
	}
	result := expandProviderDefaultTags(raw)
	if len(result) != 0 {
		t.Fatalf("expected empty tags, got %#v", result)
	}
}

func TestExpandProviderDefaultTags_Nil(t *testing.T) {
	result := expandProviderDefaultTags(nil)
	if len(result) != 0 {
		t.Fatalf("expected empty tags for nil input, got %#v", result)
	}
}

func TestExpandProviderDefaultTags_WithTags(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"tags": map[string]interface{}{
				"env":  "prod",
				"team": "cloud-railway",
			},
		},
	}
	result := expandProviderDefaultTags(raw)
	if result["env"] != "prod" {
		t.Fatalf("expected env=prod, got %v", result["env"])
	}
	if result["team"] != "cloud-railway" {
		t.Fatalf("expected team=cloud-railway, got %v", result["team"])
	}
}

func TestExpandProviderIgnoreTags_Empty(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"keys":         schema.NewSet(schema.HashString, []interface{}{}),
			"key_prefixes": schema.NewSet(schema.HashString, []interface{}{}),
		},
	}
	result := expandProviderIgnoreTags(raw)
	if len(result.Keys) != 0 {
		t.Fatalf("expected empty keys, got %#v", result.Keys)
	}
	if len(result.KeyPrefixes) != 0 {
		t.Fatalf("expected empty key_prefixes, got %#v", result.KeyPrefixes)
	}
}

func TestExpandProviderIgnoreTags_Nil(t *testing.T) {
	result := expandProviderIgnoreTags(nil)
	if len(result.Keys) != 0 {
		t.Fatalf("expected empty keys for nil input, got %#v", result.Keys)
	}
	if len(result.KeyPrefixes) != 0 {
		t.Fatalf("expected empty key_prefixes for nil input, got %#v", result.KeyPrefixes)
	}
}

func TestExpandProviderIgnoreTags_WithKeysAndPrefixes(t *testing.T) {
	raw := []interface{}{
		map[string]interface{}{
			"keys":         schema.NewSet(schema.HashString, []interface{}{"ignoreKey1", "ignoreKey2"}),
			"key_prefixes": schema.NewSet(schema.HashString, []interface{}{"aws:", "internal:"}),
		},
	}
	result := expandProviderIgnoreTags(raw)
	if _, ok := result.Keys["ignoreKey1"]; !ok {
		t.Fatalf("expected ignoreKey1 in Keys, got %#v", result.Keys)
	}
	if _, ok := result.Keys["ignoreKey2"]; !ok {
		t.Fatalf("expected ignoreKey2 in Keys, got %#v", result.Keys)
	}
	if _, ok := result.KeyPrefixes["aws:"]; !ok {
		t.Fatalf("expected aws: in KeyPrefixes, got %#v", result.KeyPrefixes)
	}
	if _, ok := result.KeyPrefixes["internal:"]; !ok {
		t.Fatalf("expected internal: in KeyPrefixes, got %#v", result.KeyPrefixes)
	}
}
