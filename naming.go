package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

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
