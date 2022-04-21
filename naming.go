package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"regexp"
)

const MaxRoleLen = 64

// NameWithPrefix returns in order:
// the name if non-empty,
// a prefix generated name if non-empty,
// or fully generated name prefixed with "terraform-".
func NameWithPrefix(name string, namePrefix string) string {
	if name != "" {
		return name
	}

	if namePrefix != "" {
		return resource.PrefixedUniqueId(namePrefix)
	}

	return resource.UniqueId()
}

func NamePrefixFromName(name string) *string {
	namePrefixIndex := len(name) - resource.UniqueIDSuffixLength

	if namePrefixIndex <= 0 {
		return nil
	}

	namePrefix := name[:namePrefixIndex]

	return &namePrefix
}

// Validate Role name based on https://docs.aws.amazon.com/IAM/latest/APIReference/API_CreateRole.html
var validRoleName = validation.All(
	validation.StringLenBetween(1, MaxRoleLen),
	validation.StringMatch(regexp.MustCompile(`^[\w+=,.@-]+$`), "must match [\\w+=,.@-]"),
)

var validRolePrefix = validation.All(
	validation.StringLenBetween(1, MaxRoleLen - resource.UniqueIDSuffixLength),
	validation.StringMatch(regexp.MustCompile(`^[\w+=,.@-]+$`), "must match [\\w+=,.@-]"),
)