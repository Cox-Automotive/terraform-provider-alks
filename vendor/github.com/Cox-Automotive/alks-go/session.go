package alks

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// SessionRequest is used to represent a new STS session request.
type SessionRequest struct {
	SessionDuration int `json:"sessionTime"`
}

// SessionResponse is used to represent a new STS session.
type SessionResponse struct {
	BaseResponse
	AccessKey       string    `json:"accessKey"`
	SecretKey       string    `json:"secretKey"`
	SessionToken    string    `json:"sessionToken"`
	SessionDuration int       `json:"sessionDuration"`
	Expires         time.Time `json:"expires"`
}

// SkypieaAccount is used to represent Skypiea data
type SkypieaAccount struct {
	Account string `json:"Account"`
	Alias   string `json:"alias"`
	Label   string `json:"label"`
}

// AccountRole is used to represent an ALKS account and role combination
type AccountRole struct {
	Account        string         `json:"account"`
	Role           string         `json:"role"`
	IamActive      bool           `json:"iamKeyActive"`
	SkypieaAccount SkypieaAccount `json:"skypieaAccount"`
}

// AccountsResponseInt is used internally to represent a collection of ALKS accounts
type AccountsResponseInt struct {
	BaseResponse
	Accounts map[string][]AccountRole `json:"accountListRole"`
}

// AccountsResponse is used to represent a collection of ALKS accounts
type AccountsResponse struct {
	Accounts []AccountRole `json:"accountListRole"`
}

// GetAccounts return a list of AccountRoles for an AWS account
func (c *Client) GetAccounts() (*AccountsResponse, *AlksError) {
	log.Printf("[INFO] Requesting available accounts from ALKS")

	b, err := json.Marshal(c.Credentials)

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error encoding account request JSON: %s", err),
		}
	}

	req, err := c.NewRequest(b, "POST", "/getAccounts/")
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

	_accts := new(AccountsResponseInt)
	err = decodeBody(resp, &_accts)
	reqID := GetRequestID(resp)
	if err != nil {

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error parsing get accounts response: %s", err),
		}
	}

	if _accts.RequestFailed() {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  _accts.BaseResponse.RequestID,
			Err:        fmt.Errorf("Error getting accounts : %s", strings.Join(_accts.GetErrors(), ", ")),
		}
	}

	accts := new(AccountsResponse)
	for k, v := range _accts.Accounts {
		v[0].Account = k
		accts.Accounts = append(accts.Accounts, v[0])
	}

	return accts, nil
}

// CreateSession will create a new STS session on AWS. If no error is
// returned then you will receive a SessionResponse object representing
// your STS session.
func (c *Client) CreateSession(sessionDuration int, useIAM bool) (*SessionResponse, *AlksError) {
	log.Printf("[INFO] Creating %v hr session", sessionDuration)

	var found = false
	durations, err := c.Durations()
	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error fetching allowable durations from ALKS: %s", err),
		}
	}

	for _, v := range durations {
		if sessionDuration == v {
			found = true
		}
	}

	if !found {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Unsupported session duration"),
		}
	}

	session := SessionRequest{sessionDuration}

	b, err := json.Marshal(struct {
		SessionRequest
		AccountDetails
	}{session, c.AccountDetails})

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error encoding session create JSON: %s", err),
		}
	}

	var endpoint = "/getKeys/"
	if useIAM {
		endpoint = "/getIAMKeys/"
	}
	req, err := c.NewRequest(b, "POST", endpoint)
	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}

	resp, httpErr := c.http.Do(req)
	if httpErr != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  "",
			Err:        err,
		}
	}

	sr := new(SessionResponse)
	err = decodeBody(resp, &sr)

	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf("Error parsing session create response: %s", err),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  "",
			Err:        fmt.Errorf("Error parsing session create response: %s", err),
		}
	}

	if sr.RequestFailed() {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  sr.BaseResponse.RequestID,
			Err:        fmt.Errorf("Error creating session: %s", strings.Join(sr.GetErrors(), ", ")),
		}
	}

	sr.Expires = time.Now().Local().Add(time.Hour * time.Duration(sessionDuration))
	sr.SessionDuration = sessionDuration

	return sr, nil
}
