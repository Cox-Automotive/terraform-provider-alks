package alks

import (
	"errors"
	"regexp"
	"strings"
)

// This regex will attempt to parse ALKS account strings (and must be valid for the package to compile)
// *** WARNING ***: The group names in the regex are referenced below, changing them means updating the associated methods as well
var accountRegex = regexp.MustCompile(`(?P<AccountNumber>\d+)(/(?P<RoleName>(ALKS)?\w+)(\s-\s(?P<AccountDesc>\w+))?)?`)

// AccountDetails represents the callers Account and Role information for ALKS requests
type AccountDetails struct {
	Account string `json:"account,omitempty"`
	Role    string `json:"role,omitempty"`
}

// GetAccountNumber parses the Account provided in AccountDetails and returns the account number if present
func (a AccountDetails) GetAccountNumber() (string, error) {
	if a.Account == "" {
		return "", errors.New("Account is empty")
	}

	if accountRegex.MatchString(a.Account) {
		matches := accountRegex.FindStringSubmatch(a.Account)

		for i, v := range accountRegex.SubexpNames() {
			if v == "AccountNumber" {
				return matches[i], nil
			}
		}
	}

	return "", errors.New("Invalid Account format")
}

// GetRoleName returns the AccountDetails Role or parses the role value from the Account
func (a AccountDetails) GetRoleName(stripPrefix bool) (string, error) {
	if a.Role != "" {
		if stripPrefix {
			return strings.TrimPrefix(a.Role, "ALKS"), nil
		}

		return a.Role, nil
	}

	if a.Account == "" {
		return "", errors.New("Account is empty")
	}

	if accountRegex.MatchString(a.Account) {
		matches := accountRegex.FindStringSubmatch(a.Account)

		for i, v := range accountRegex.SubexpNames() {
			if v == "RoleName" {
				roleName := matches[i]
				if roleName == "" {
					return "", errors.New("No Role found")
				}

				if stripPrefix {
					return strings.TrimPrefix(roleName, "ALKS"), nil
				}

				return roleName, nil
			}
		}
	}

	return "", errors.New("Invalid Account format")
}

// GetAccountDesc parses the Account provided in AccountDetails and returns the account description if present
func (a AccountDetails) GetAccountDesc() (string, error) {
	if a.Account == "" {
		return "", errors.New("Account is empty")
	}

	if accountRegex.MatchString(a.Account) {
		matches := accountRegex.FindStringSubmatch(a.Account)

		for i, v := range accountRegex.SubexpNames() {
			if v == "AccountDesc" {
				if matches[i] == "" {
					return "", errors.New("No AccountDesc found")
				}

				return matches[i], nil
			}
		}
	}

	return "", errors.New("Invalid Account format")
}
