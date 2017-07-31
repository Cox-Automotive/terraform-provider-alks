package main

import (
	"fmt"
	"github.com/Cox-Automotive/alks-go"
	"log"
)

type Config struct {
	Url      string
	Username string
	Password string
	Account  string
	Role     string
}

func (c *Config) Client() (*alks.Client, error) {

	testClient, testErr := alks.NewClient(c.Url, c.Username, c.Password, "", "")

	if testErr != nil {
		return nil, testErr
	}

	// retreive list of accounts the user has access to
	resp, err := testClient.GetAccounts()

	if err != nil {
		// this error check will catch invalid user credentials
		return nil, err
	} else {
		// this check will catch incorrect accounts/role
		var validAccountRole = false
		for _, acct := range resp.Accounts {
			if acct.Account == c.Account && acct.Role == c.Role {
				validAccountRole = true
			}
		}

		if !validAccountRole {
			return nil, fmt.Errorf("The specified account (%v) and role (%v) are not available in your ALKS account. Please verify your permissions and try again.", c.Account, c.Role)
		}
	}

	// now we know we have a valid user credentials and account/role so create the client
	client, err := alks.NewClient(c.Url, c.Username, c.Password, c.Account, c.Role)

	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] ALKS Client configured")

	return client, nil
}
