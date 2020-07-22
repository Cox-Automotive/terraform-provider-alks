package alks

import (
	"fmt"
	"log"
	"strings"
)

// GetMyLoginRole returns the LoginRole corresponding to the clients current STS credentials
func (c *Client) GetMyLoginRole() (*LoginRoleResponse, error) {
	log.Printf("[INFO] Requesting Login Role information from ALKS")

	if !c.IsUsingSTSCredentials() {
		return nil, fmt.Errorf("GetMyLoginRole only supports clients using STS credentials, try using GetLoginRole instead")
	}

	req, err := c.NewRequest(nil, "GET", "/loginRoles/id/me")
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	lrr := new(LoginRoleResponse)
	err = decodeBody(resp, &lrr)
	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf("Error parsing LoginRole response: [%s] %s", reqID, err)
		}

		return nil, fmt.Errorf("Error parsing LoginRole response: %s", err)
	}

	if lrr.RequestFailed() {
		return nil, fmt.Errorf("Error fetching role information: [%s] %s", lrr.BaseResponse.RequestID, strings.Join(lrr.GetErrors(), ", "))
	}

	return lrr, nil
}

// GetLoginRole returns the login role corresponding to the current account and role stored in AccountDetails
func (c *Client) GetLoginRole() (*LoginRoleResponse, error) {
	// If the client is configured with STS call the correct method
	if c.IsUsingSTSCredentials() {
		log.Println("[INFO] Client configured with STS credentials, dispatching to GetMyLoginRole instead")
		return c.GetMyLoginRole()
	}

	account, err := c.AccountDetails.GetAccountNumber()
	if err != nil {
		return nil, err
	}

	roleName, err := c.AccountDetails.GetRoleName(false)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Requesting Login Role information for %v/%v from ALKS", account, roleName)

	req, err := c.NewRequest(nil, "GET", fmt.Sprintf("/loginRoles/id/%v/%v", account, roleName))
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	lrr := new(LoginRoleResponse)
	err = decodeBody(resp, &lrr)
	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf("Error parsing LoginRole response: [%s] %s", reqID, err)
		}

		return nil, fmt.Errorf("Error parsing LoginRole response: %s", err)
	}

	if lrr.RequestFailed() {
		return nil, fmt.Errorf("Error fetching role information: [%s] %s", lrr.BaseResponse.RequestID, strings.Join(lrr.GetErrors(), ", "))
	}

	return lrr, nil
}
