package main

import (
	"fmt"

	"github.com/Cox-Automotive/alks-go"
)

func validateIAMEnabled(client *alks.Client) *alks.AlksError {
	// Validate STS for IAM active.
	resp, err := client.GetMyLoginRole()
	if err != nil {
		return err
	}

	if !resp.LoginRole.IamKeyActive {
		return &alks.AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err: fmt.Errorf("uh oh! You're using the " + resp.LoginRole.Role + " role which is not IAM active. " +
				"Please instead use one of the following roles: Admin, IAMAdmin, LabAdmin, or a Machine Identity"),
		}
	}

	return nil
}
