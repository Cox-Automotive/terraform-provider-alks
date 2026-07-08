package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestMigrateV0toV1_AddsEnableAlksAccess(t *testing.T) {
	state := &terraform.InstanceState{
		ID: "test-role-id",
		Attributes: map[string]string{
			"name": "test-role",
			"type": "Amazon EC2",
		},
	}
	result, err := migrateV0toV1(state)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if result.Attributes["enable_alks_access"] != "false" {
		t.Fatalf("expected enable_alks_access=false, got %s", result.Attributes["enable_alks_access"])
	}
}

func TestMigrateV0toV1_EmptyState(t *testing.T) {
	state := &terraform.InstanceState{}
	result, err := migrateV0toV1(state)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if _, ok := result.Attributes["enable_alks_access"]; ok {
		t.Fatal("enable_alks_access should not be set on an empty InstanceState")
	}
}

func TestMigrateV0toV1_PreservesExistingValue(t *testing.T) {
	state := &terraform.InstanceState{
		ID: "test-role-id",
		Attributes: map[string]string{
			"name":               "test-role",
			"enable_alks_access": "true",
		},
	}
	result, err := migrateV0toV1(state)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if result.Attributes["enable_alks_access"] != "true" {
		t.Fatalf("expected existing enable_alks_access=true to be preserved, got %s", result.Attributes["enable_alks_access"])
	}
}

func TestMigrateV0toV1_PreservesOtherAttributes(t *testing.T) {
	state := &terraform.InstanceState{
		ID: "test-role-id",
		Attributes: map[string]string{
			"name": "my-role",
			"type": "Amazon EC2",
			"arn":  "arn:aws:iam::123456789012:role/my-role",
		},
	}
	result, err := migrateV0toV1(state)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if result.Attributes["name"] != "my-role" {
		t.Fatalf("expected name=my-role to be preserved, got %s", result.Attributes["name"])
	}
	if result.Attributes["arn"] != "arn:aws:iam::123456789012:role/my-role" {
		t.Fatalf("arn should be preserved, got %s", result.Attributes["arn"])
	}
}

func TestMigrateState_V0ToV1(t *testing.T) {
	state := &terraform.InstanceState{
		ID:         "test-id",
		Attributes: map[string]string{"name": "test"},
	}
	result, err := migrateState(0, state, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if result.Attributes["enable_alks_access"] != "false" {
		t.Fatalf("expected false after v0->v1 migration, got %s", result.Attributes["enable_alks_access"])
	}
}

func TestMigrateState_UnknownVersion(t *testing.T) {
	state := &terraform.InstanceState{
		ID:         "test-id",
		Attributes: map[string]string{"name": "test"},
	}
	_, err := migrateState(99, state, nil)
	if err == nil {
		t.Fatal("expected error for unknown schema version")
	}
}
