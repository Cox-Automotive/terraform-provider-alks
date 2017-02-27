package alks

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// SessionRequest is used to represent a new STS session request.
type SessionRequest struct {
	SessionDuration int `json:"sessionTime"`
}

// SessionResponse is used to represent a new STS session.
type SessionResponse struct {
	AccessKey       string    `json:"accessKey"`
	SecretKey       string    `json:"secretKey"`
	SessionToken    string    `json:"sessionToken"`
	SessionDuration int       `json:"sessionDuration"`
	Expires         time.Time `json:"expires"`
}

// AccountRole is used to represent an ALKS account and role combination
type AccountRole struct {
	Account   string `json:"account"`
	Role      string `json:"role"`
	IamActive bool   `json:"iamKeyActive"`
}

// AccountsResponseInt is used internally to represent a collection of ALKS accounts
type AccountsResponseInt struct {
	Accounts map[string][]AccountRole `json:"accountListRole"`
}

// AccountsResponse is used to represent a collection of ALKS accounts
type AccountsResponse struct {
	Accounts []AccountRole `json:"accountListRole"`
}

func (c *Client) GetAccounts() (*AccountsResponse, error) {
	log.Printf("[INFO] Requesting available accounts from ALKS")

	b, err := json.Marshal(c.Account)

	if err != nil {
		return nil, fmt.Errorf("Error encoding account request JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/getAccounts/")
	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return nil, err
	}

	_accts := new(AccountsResponseInt)
	err = decodeBody(resp, &_accts)

	if err != nil {
		return nil, fmt.Errorf("Error parsing session create response: %s", err)
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
func (c *Client) CreateSession(sessionDuration int) (*SessionResponse, error) {
	log.Printf("[INFO] Creating %v hr session", sessionDuration)

	var found bool = false
	for _, v := range c.Durations() {
		if sessionDuration == v {
			found = true
		}
	}

	if !found {
		return nil, fmt.Errorf("Unsupported session duration")
	}

	session := SessionRequest{sessionDuration}

	b, err := json.Marshal(struct {
		SessionRequest
		AlksAccount
	}{session, c.Account})

	if err != nil {
		return nil, fmt.Errorf("Error encoding session create JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/getKeys/")
	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return nil, err
	}

	sr := new(SessionResponse)
	err = decodeBody(resp, &sr)

	if err != nil {
		return nil, fmt.Errorf("Error parsing session create response: %s", err)
	}

	sr.Expires = time.Now().Local().Add(time.Hour * time.Duration(sessionDuration))
	sr.SessionDuration = sessionDuration

	return sr, nil
}
