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
		return errors.New("Looks like you are not using ALKS IAM credentials. This will result in errors when creating roles. \n " +
			"Note: If using ALKS CLI to get credentials, be sure to use the '-i' flag. \n Please see https://coxautoinc.sharepoint.com/sites/service-internal-tools-team/SitePages/ALKS-Terraform-Provider---Troubleshooting.aspx for more information.")
	}

	return nil
}
