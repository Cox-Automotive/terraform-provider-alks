package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	awspolicy "github.com/hashicorp/awspolicyequivalence"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type TrustPolicyDocument struct {
	Version    string      `json:",omitempty"`
	Id         string      `json:",omitempty"`
	Statements interface{} `json:"Statement"`
}

func SuppressEquivalentTrustPolicyDiffs(key string, old string, new string, d *schema.ResourceData) bool {
	if strings.TrimSpace(old) == "" && strings.TrimSpace(new) == "" {
		return true
	}

	if strings.TrimSpace(old) == "{}" && strings.TrimSpace(new) == "" {
		return true
	}

	if strings.TrimSpace(old) == "" && strings.TrimSpace(new) == "{}" {
		return true
	}

	if strings.TrimSpace(old) == "{}" && strings.TrimSpace(new) == "{}" {
		return true
	}

	equivalent, err := awspolicy.PoliciesAreEquivalent(old, new)
	if err != nil {
		return false
	}

	return equivalent
}

// Broken into seperate function to allow for returning of errors.
func UnmarshalAndMarshal(policy []byte) ([]byte, error) {
	unmarshaledPolicy := TrustPolicyDocument{}
	if err := json.Unmarshal(policy, &unmarshaledPolicy); err != nil {
		return nil, fmt.Errorf("Error unmarshaling trust policy: %s", err)
	}

	marshaledPolicy, err := json.Marshal(unmarshaledPolicy)
	if err != nil {
		return nil, fmt.Errorf("Error marshaling trust policy: %s", err)
	}

	return marshaledPolicy, nil
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
