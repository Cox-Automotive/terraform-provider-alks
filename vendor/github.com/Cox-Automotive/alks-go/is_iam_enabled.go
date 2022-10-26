package alks

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type IsIamEnabledRequest struct {
	AccountDetails
	RoleArn string `json:"roleArn,omitempty"`
}

// IsIamEnabledResponse is used to represent a role that's IAM active or not.
type IsIamEnabledResponse struct {
	BaseResponse
	AccountDetails
	RoleArn    string `json:"roleArn"`
	IamEnabled bool   `json:"iamEnabled"`
}

// IsIamEnabled will check if a MI, AccountDetails, or STS assumed role is IAM active or not.
func (c *Client) IsIamEnabled(roleArn string) (*IsIamEnabledResponse, *AlksError) {

	if len(roleArn) > 1 {
		log.Printf("[INFO] Is IAM enabled for MI: %s", roleArn)
	} else {
		log.Printf("[INFO] Is IAM enabled for: %s/%s", c.AccountDetails.Account, c.AccountDetails.Role)
	}

	iam := IsIamEnabledRequest{
		c.AccountDetails,
		roleArn,
	}

	body, err := json.Marshal(iam)

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("error encoding IAM create role JSON: %s", err),
		}

	}

	req, err := c.NewRequest(body, "POST", "/isIamEnabled")
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

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		iamErr := new(AlksResponseError)
		err = decodeBody(resp, &iamErr)
		if err != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, &AlksError{
					StatusCode: resp.StatusCode,
					RequestId:  reqID,
					Err:        fmt.Errorf(ParseErrorReqId, reqID, err),
				}
			}

			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  "",
				Err:        fmt.Errorf(ParseError, err),
			}
		}

		if iamErr.Errors != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, &AlksError{
					StatusCode: resp.StatusCode,
					RequestId:  reqID,
					Err:        fmt.Errorf(ErrorStringFull, reqID, resp.StatusCode, iamErr.Errors),
				}
			}

			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  "",
				Err:        fmt.Errorf(ErrorStringNoReqId, resp.StatusCode, iamErr.Errors),
			}
		}

		if reqID := GetRequestID(resp); reqID != "" {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ErrorStringOnlyCodeAndReqId, reqID, resp.StatusCode),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  "",
			Err:        fmt.Errorf(ErrorStringOnlyCode, resp.StatusCode),
		}
	}

	validate := new(IsIamEnabledResponse)
	err = decodeBody(resp, validate)

	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf("error parsing isIamEnabled response: [%s] %s", reqID, err),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  "",
			Err:        fmt.Errorf("error parsing isIamEnabled response: %s", err),
		}
	}
	if validate.RequestFailed() {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  validate.BaseResponse.RequestID,
			Err:        fmt.Errorf("error validating if IAM enabled: [%s] %s", validate.BaseResponse.RequestID, strings.Join(validate.GetErrors(), ", ")),
		}
	}

	return validate, nil
}
