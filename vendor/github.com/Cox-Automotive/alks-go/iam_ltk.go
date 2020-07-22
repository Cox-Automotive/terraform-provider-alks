package alks

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

// LongTermKey represents a long term key
type LongTermKey struct {
	UserName    string `json:"userName"`
	AccessKeyID string `json:"accessKeyId"`
	Status      string `json:"status"`
	CreateDate  string `json:"createDate"`
}

// GetLongTermKeysResponse is used to represent the list of long term keys
type GetLongTermKeysResponse struct {
	BaseResponse
	LongTermKeys []LongTermKey `json:"longTermKeys"`
}

// GetLongTermKeyResponse is used to represent a single long term key.
type GetLongTermKeyResponse struct {
	BaseResponse
	LongTermKey `json:"longTermKey"`
}

// BaseLongTermKeyResponse encapsulates shared response fields
type BaseLongTermKeyResponse struct {
	Action              string `json:"action,omitempty"`
	AddedIAMUserToGroup bool   `json:"addedIAMUserToGroup,omitempty"`
	PartialError        bool   `json:"partialError,omitempty"`
}

// CreateLongTermKey represents the response from API
type CreateLongTermKey struct {
	IAMUserName string `json:"iamUserName"`
	IAMUserArn  string `json:"iamUserArn"`
	AccessKey   string `json:"accessKey"`
	SecretKey   string `json:"secretKey"`
}

// LongTermKeyRequest is used to represent the request body to create or delete LTKs
type LongTermKeyRequest struct {
	AccountDetails
	IamUserName string `json:"iamUserName"`
}

// CreateLongTermKeyResponse is the response to the CLI client
type CreateLongTermKeyResponse struct {
	AccountDetails
	BaseResponse
	BaseLongTermKeyResponse
	CreateLongTermKey
}

// DeleteLongTermKeyResponse is the response to the CLI client
type DeleteLongTermKeyResponse struct {
	AccountDetails
	BaseResponse
	BaseLongTermKeyResponse
}

// GetLongTermKeys gets the LTKs for an account
// If no error is returned then you will receive a list of LTKs
func (c *Client) GetLongTermKeys() (*GetLongTermKeysResponse, error) {
	log.Printf("[INFO] Getting long term keys")

	accountID, err := c.AccountDetails.GetAccountNumber()
	if err != nil {
		return nil, fmt.Errorf("Error reading Account value: %s", err)
	}

	roleName, err := c.AccountDetails.GetRoleName(false)
	if err != nil {
		return nil, fmt.Errorf("Error reading Role value: %s", err)
	}

	req, err := c.NewRequest(nil, "GET", "/ltks/"+accountID+"/"+roleName)
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	cr := new(GetLongTermKeysResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf("Error parsing GetLongTermKeysResponse: [%s] %s", reqID, err)
		}

		return nil, fmt.Errorf("Error parsing GetLongTermKeysResponse: %s", err)
	}

	if cr.RequestFailed() {
		return nil, fmt.Errorf("Error getting long term keys: [%s] %s", cr.BaseResponse.RequestID, strings.Join(cr.GetErrors(), ", "))
	}

	return cr, nil
}

// GetLongTermKey gets a single LTK for an account
// If no error is returned, then you will receive an LTK for the given account.
func (c *Client) GetLongTermKey(iamUsername string) (*GetLongTermKeyResponse, error) {
	log.Printf("[INFO] Getting long term key")

	var req *http.Request
	var err error

	if c.IsUsingSTSCredentials() {
		req, err = c.NewRequest(nil, "GET", "/ltk/search/"+iamUsername)
	} else {
		accountID, err := c.AccountDetails.GetAccountNumber()
		if err != nil {
			return nil, fmt.Errorf("error reading Account value: %s", err)
		}

		roleName, err := c.AccountDetails.GetRoleName(false)
		if err != nil {
			return nil, fmt.Errorf("error reading Role value: %s", err)
		}

		req, err = c.NewRequest(nil, "GET", "/ltk/"+accountID+"/"+roleName+"/search/"+iamUsername)
	}

	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	cr := new(GetLongTermKeyResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf("error parsing GetLongTermKeyResponse: [%s] %s", reqID, err)
		}
	}

	if cr.RequestFailed() {
		return nil, fmt.Errorf("error getting long term keys: [%s] %s", cr.BaseResponse.RequestID, strings.Join(cr.GetErrors(), ", "))
	}

	return cr, nil
}

// CreateLongTermKey creates an LTK user for an account.
// If no error is returned, then you will receive an appropriate success message.
func (c *Client) CreateLongTermKey(iamUsername string) (*CreateLongTermKeyResponse, error) {
	log.Printf("[INFO] Creating long term key: %s", iamUsername)

	request := LongTermKeyRequest{
		AccountDetails: c.AccountDetails,
		IamUserName:    iamUsername,
	}

	reqBody, err := json.Marshal(request)

	if err != nil {
		return nil, fmt.Errorf("error encoding LTK create JSON: %s", err)
	}

	req, err := c.NewRequest(reqBody, "POST", "/accessKeys")
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	cr := new(CreateLongTermKeyResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf("error parsing CreateLongTermKeyResponse: [%s] %s", reqID, err)
		}

		return nil, fmt.Errorf("error parsing CreateLongTermKeyResponse: %s", err)
	}

	if cr.RequestFailed() {
		return nil, fmt.Errorf("error creating long term key: [%s] %s", cr.BaseResponse.RequestID, strings.Join(cr.GetErrors(), ", "))
	}

	return cr, nil
}

// DeleteLongTermKey deletes an LTK user for an account.
// If no error is returned, then you will receive an appropriate success message.
func (c *Client) DeleteLongTermKey(iamUsername string) (*DeleteLongTermKeyResponse, error) {
	log.Printf("[INFO] Deleting long term key: %s", iamUsername)

	request := LongTermKeyRequest{
		AccountDetails: c.AccountDetails,
		IamUserName:    iamUsername,
	}

	reqBody, err := json.Marshal(request)

	if err != nil {
		return nil, fmt.Errorf("error encoding LTK delete JSON: %s", err)
	}

	req, err := c.NewRequest(reqBody, "DELETE", "/IAMUser")
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	cr := new(DeleteLongTermKeyResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf("error parsing DeleteLongTermKeyResponse: [%s] %s", reqID, err)
		}

		return nil, fmt.Errorf("error parsing DeleteLongTermKeyResponse: %s", err)
	}

	if cr.RequestFailed() {
		return nil, fmt.Errorf("error deleting long term key: [%s] %s", cr.BaseResponse.RequestID, strings.Join(cr.GetErrors(), ", "))
	}

	return cr, nil

}
