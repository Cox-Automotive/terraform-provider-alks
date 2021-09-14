package main

import (
	"errors"

	"github.com/Cox-Automotive/alks-go"
)

func validateIAMEnabled(client *alks.Client) error {
	// Validate STS for IAM active.
	resp, err := client.GetMyLoginRole()
	if err != nil {
		return err
	}

	if !resp.LoginRole.IamKeyActive {
		return errors.New("uh oh! You're using the " + resp.LoginRole.Role + " which is not IAM active. " +
			"Please instead use one of the following roles: Admin, IAMAdmin, LabAdmin, or a Machine Identity")
	}

	return nil
}
