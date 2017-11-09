package alks

import (
	"encoding/json"
	"fmt"
	"log"
)

// IamSessionRequest is used to represent a new IAM session request.
type IamSessionRequest struct {
	SessionTime int `json:"sessionTime"`
}

// IamSessionResponse is used to represent the session that is created.
type IamSessionResponse struct {
	AccessKey    string `json:"accessKey"`
	SessionKey   string `json:"secretKey"`
	SessionToken string `json:"sessionToken"`
}

// CreateIamSession creates a new IAM STS session. If no error is returned
// then you will received a IamSessionResponse object containing your session
// keys.
func (c *Client) CreateIamSession() (*IamSessionResponse, error) {
	log.Println("[INFO] Creating IAM session")

	iam := IamSessionRequest{1}
	b, err := json.Marshal(struct {
		IamSessionRequest
		AlksAccount
	}{iam, c.Account})

	if err != nil {
		return nil, fmt.Errorf("Error encoding IAM new session JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/getIAMKeys/")
	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return nil, err
	}

	ses := new(IamSessionResponse)
	err = decodeBody(resp, &ses)

	if err != nil {
		return nil, fmt.Errorf("Error parsing create IAM session response: %s", err)
	}

	return ses, nil
}
