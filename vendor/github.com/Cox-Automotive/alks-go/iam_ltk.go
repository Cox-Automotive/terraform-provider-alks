package alks

import (
	"fmt"
	"log"
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

// GetLongTermKeys gets the LTKs for an account
// If no error is returned then you will receive a list of LTKs
func (c *Client) GetLongTermKeys(accountId string, roleName string) (*GetLongTermKeysResponse, error) {
	log.Printf("[INFO] Getting long term keys for: %s/%s", accountId, roleName)

	req, err := c.NewRequest(nil, "GET", "/ltks/"+accountId+"/"+roleName)
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
		if reqId := GetRequestID(resp); reqId != "" {
			return nil, fmt.Errorf("Error parsing GetLongTermKeysResponse: [%s] %s", reqId, err)
		}

		return nil, fmt.Errorf("Error parsing GetLongTermKeysResponse: %s", err)
	}

	if cr.RequestFailed() {
		return nil, fmt.Errorf("Error getting long term keys: [%s] %s", cr.BaseResponse.RequestID, strings.Join(cr.GetErrors(), ", "))
	}

	return cr, nil
}
