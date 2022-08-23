package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TrustPolicyDocument struct {
	Version    string      `json:",omitempty"`
	Id         string      `json:",omitempty"`
	Statements interface{} `json:"Statement"`
}

func SuppressEquivalentTrustPolicyDiffs(key string, old string, new string, d *schema.ResourceData) bool {
	UnmarshalAndMarshal(&old, &new)
	return old == new
}

//Broken into seperate function to allow for returning of errors.
func UnmarshalAndMarshal(oldPolicy *string, newPolicy *string) error {
	policyOldIntermediate := TrustPolicyDocument{}
	policyNewIntermediate := TrustPolicyDocument{}
	if err := json.Unmarshal([]byte(*oldPolicy), &policyOldIntermediate); err != nil {
		return fmt.Errorf("Error unmarshaling old trust policy: %s", err)
	}
	if err := json.Unmarshal([]byte(*newPolicy), &policyNewIntermediate); err != nil {
		return fmt.Errorf("Error unmarshaling new trust policy: %s", err)
	}

	oldByteArray, err := json.Marshal(policyOldIntermediate)
	if err != nil {
		return fmt.Errorf("Error marshaling old trust policy: %s", err)
	}
	newByteArray, err := json.Marshal(policyNewIntermediate)
	if err != nil {
		return fmt.Errorf("Error marshaling new trust policy: %s", err)
	}
	*oldPolicy = string(oldByteArray)
	*newPolicy = string(newByteArray)
	return nil
}

// Using a diff function is the currently accepted way to compare the configuration of two different attributes at plan time.
func trustPoliciesWithIncludeDefaultPolicies(_ context.Context, diff *schema.ResourceDiff, meta interface{}) error {
	var assumeRolePolicy = diff.Get("assume_role_policy")
	var includeDefaultPolicies = (diff.Get("include_default_policies").(bool))
	if (assumeRolePolicy != nil) && (assumeRolePolicy != "") {
		if includeDefaultPolicies {
			return fmt.Errorf("include_default_policies must be false or excluded if including an assume_role_policy %#v", assumeRolePolicy)
		}
	}
	return nil
}
