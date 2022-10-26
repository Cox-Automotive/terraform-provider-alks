package alks

import (
	"fmt"
	"log"
	"strings"
)

// GetMyLoginRole returns the LoginRole corresponding to the clients current STS credentials
func (c *Client) GetMyLoginRole() (*LoginRoleResponse, *AlksError) {
	log.Printf("[INFO] Requesting Login Role information from ALKS")

	if !c.IsUsingSTSCredentials() {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("GetMyLoginRole only supports clients using STS credentials, try using GetLoginRole instead"),
		}
	}

	req, err := c.NewRequest(nil, "GET", "/loginRoles/id/me")
	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}

	reqID := GetRequestID(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		loginErr := new(AlksResponseError)
		err = decodeBody(resp, &loginErr)
		if err != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseError, err),
			}
		}

		if loginErr.Errors != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(AlksResponsErrorStrings, loginErr.Errors),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}

	}

	lrr := new(LoginRoleResponse)
	err = decodeBody(resp, &lrr)
	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, &AlksError{
				StatusCode: 0,
				RequestId:  reqID,
				Err:        fmt.Errorf("Error parsing LoginRole response: %s", err),
			}
		}

		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error parsing LoginRole response: %s", err),
		}
	}

	if lrr.RequestFailed() {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  lrr.BaseResponse.RequestID,
			Err:        fmt.Errorf("Error fetching role information: %s", strings.Join(lrr.GetErrors(), ", ")),
		}
	}

	return lrr, nil
}

// GetLoginRole returns the login role corresponding to the current account and role stored in AccountDetails
func (c *Client) GetLoginRole() (*LoginRoleResponse, *AlksError) {
	// If the client is configured with STS call the correct method
	if c.IsUsingSTSCredentials() {
		log.Println("[INFO] Client configured with STS credentials, dispatching to GetMyLoginRole instead")
		return c.GetMyLoginRole()
	}

	account, err := c.AccountDetails.GetAccountNumber()
	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}

	roleName, err := c.AccountDetails.GetRoleName(false)
	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}

	log.Printf("[INFO] Requesting Login Role information for %v/%v from ALKS", account, roleName)

	req, err := c.NewRequest(nil, "GET", fmt.Sprintf("/loginRoles/id/%v/%v", account, roleName))
	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  "",
			Err:        err,
		}
	}

	reqID := GetRequestID(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		loginErr := new(AlksResponseError)
		err = decodeBody(resp, &loginErr)
		if err != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseError, err),
			}
		}

		if loginErr.Errors != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(AlksResponsErrorStrings, loginErr.Errors),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}
	}

	lrr := new(LoginRoleResponse)
	err = decodeBody(resp, &lrr)
	if err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error parsing LoginRole response: %s", err),
		}
	}

	if lrr.RequestFailed() {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  lrr.BaseResponse.RequestID,
			Err:        fmt.Errorf("Error fetching role information: %s", strings.Join(lrr.GetErrors(), ", ")),
		}
	}

	return lrr, nil
}
