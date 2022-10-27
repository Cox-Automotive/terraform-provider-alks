package alks

import (
	"encoding/json"
	"fmt"
	"log"

	// "net/http"
	"strings"
)

//Represents iamUser returned by iam-user endpoint
type IamUser struct {
	ARN       string `json:"arn"`
	AccountId string `json:"accountId"`
	UserName  string `json:"userName"`
	AccessKey string `json:"accessKey"`
	Tags      []Tag  `json:"tags"`
}

// AllIamUsersResponseType represents iamUser returned by ltks endpoint
type AllIamUsersResponseType struct {
	UserName    string `json:"userName"`
	AccessKeyID string `json:"accessKeyId"`
	Status      string `json:"status"`
	CreateDate  string `json:"createDate"`
}

// GetIamUsersResponse is used to represent the list of long term keys
type GetIamUsersResponse struct {
	BaseResponse
	IamUsers []AllIamUsersResponseType `json:"longTermKeys"`
}

// GetIamUserResponse is used to represent a single long term key.
type GetIamUserResponse struct {
	BaseResponse
	User IamUser `json:"item"`
}

// BaseIamUserResponse encapsulates shared response fields
type BaseIamUserResponse struct {
	Action              string `json:"action,omitempty"`
	AddedIAMUserToGroup bool   `json:"addedIAMUserToGroup,omitempty"`
	PartialError        bool   `json:"partialError,omitempty"`
}

// CreateIamUserApiResponse represents the response from API
type CreateIamUserApiResponse struct {
	IAMUserName string `json:"iamUserName"`
	IAMUserArn  string `json:"iamUserArn"`
	AccessKey   string `json:"accessKey"`
	SecretKey   string `json:"secretKey"`
}

type CreateIamUserRequest struct {
	AccountDetails
	IamUserName string `json:"iamUserName"`
	Tags        []Tag  `json:"tags,omitempty"`
}

// CreateIamUserResponse is the response to the CLI client
type CreateIamUserResponse struct {
	AccountDetails
	BaseResponse
	BaseIamUserResponse
	CreateIamUserApiResponse
}

//Used as options for create and update iamUser
type IamUserOptions struct {
	IamUserName *string
	Tags        *[]Tag
}

// DeleteIamUserRequest is used to represent the request body to delete LTKs
type DeleteIamUserRequest struct {
	AccountDetails
	IamUserName string `json:"iamUserName"`
}

type DeleteIamUserResponse struct {
	AccountDetails
	BaseResponse
	BaseIamUserResponse
}

type UpdateIamUserRequest struct {
	User struct {
		Tags []Tag `json:"tags"`
	} `json:"user"`
}

type UpdateIamUserResponse struct {
	BaseResponse
	User IamUser `json:"item"`
}

// GetIamUsers gets the LTKs for an account
// If no error is returned then you will receive a list of LTKs
func (c *Client) GetIamUsers() (*GetIamUsersResponse, *AlksError) {
	log.Printf("[INFO] Getting long term keys")

	accountID, err := c.AccountDetails.GetAccountNumber()
	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error reading Account value: %s", err),
		}
	}

	roleName, err := c.AccountDetails.GetRoleName(false)
	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error reading Role value: %s", err),
		}
	}

	req, err := c.NewRequest(nil, "GET", "/ltks/"+accountID+"/"+roleName)
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
		keyErr := new(AlksResponseError)
		err = decodeBody(resp, &keyErr)
		if err != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseError, err),
			}
		}

		if keyErr.Errors != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(AlksResponsErrorStrings, strings.Join(keyErr.Errors, ", ")),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}
	}

	cr := new(GetIamUsersResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error parsing GetLongTermKeysResponse: %s", err),
		}
	}

	return cr, nil
}

// GetIamUser gets a single LTK for an account
// If no error is returned, then you will receive an LTK for the given account.
func (c *Client) GetIamUser(iamUsername string) (*GetIamUserResponse, *AlksError) {
	log.Printf("[INFO] Getting long term key")

	accountID, err := c.AccountDetails.GetAccountNumber()
	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error reading Account value: %s", err),
		}
	}

	req, err := c.NewRequest(nil, "GET", "/iam-users/id/"+accountID+"/"+iamUsername)

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error creating request object: %s", err),
		}
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error during request: %s", err),
		}
	}

	reqID := GetRequestID(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		keyErr := new(AlksResponseError)
		err = decodeBody(resp, &keyErr)
		if err != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseError, err),
			}
		}

		if keyErr.Errors != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, &AlksError{
					StatusCode: resp.StatusCode,
					RequestId:  reqID,
					Err:        fmt.Errorf(AlksResponsErrorStrings, strings.Join(keyErr.Errors, ", ")),
				}
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}
	}

	cr := new(GetIamUserResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("error parsing GetLongTermKeyResponse: %s", err),
		}
	}

	return cr, nil
}

func NewCreateIamUserRequest(options *IamUserOptions) (*CreateIamUserRequest, error) {
	if options.IamUserName == nil {
		return nil, fmt.Errorf("IamUserName option must not be nil")
	}

	iamUser := &CreateIamUserRequest{}
	iamUser.IamUserName = *options.IamUserName

	if options.Tags != nil {
		iamUser.Tags = *options.Tags
	} else {
		iamUser.Tags = nil
	}

	return iamUser, nil
}

// CreateIamUser creates an iamUser and secret key for an account.
// If no error is returned, then you will receive an appropriate success message.
func (c *Client) CreateIamUser(options *IamUserOptions) (*CreateIamUserResponse, *AlksError) {
	request, err := NewCreateIamUserRequest(options)

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}
	log.Printf("[INFO] Creating long term key: %s", *options.IamUserName)

	request.AccountDetails = c.AccountDetails

	log.Printf("[INFO] The request body is %v", *request)

	b, err := json.Marshal(struct {
		CreateIamUserRequest
	}{*request})

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("error encoding LTK create JSON: %s", err),
		}
	}

	log.Printf("[INFO] Request Body: %v", string(b))

	req, err := c.NewRequest(b, "POST", "/accessKeys")

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
		keyErr := new(AlksResponseError)
		err = decodeBody(resp, &keyErr)

		if err != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseError, err),
			}
		}

		if keyErr.Errors != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(AlksResponsErrorStrings, strings.Join(keyErr.Errors, ", ")),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}
	}

	cr := new(CreateIamUserResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("error parsing CreateLongTermKeyResponse: %s", err),
		}
	}
	return cr, nil
}

// DeleteIamUser deletes an LTK user for an account.
// If no error is returned, then you will receive an appropriate success message.
func (c *Client) DeleteIamUser(iamUsername string) (*DeleteIamUserResponse, *AlksError) {
	log.Printf("[INFO] Deleting long term key: %s", iamUsername)

	request := DeleteIamUserRequest{
		AccountDetails: c.AccountDetails,
		IamUserName:    iamUsername,
	}

	reqBody, err := json.Marshal(request)

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("error encoding iamUser delete JSON: %s", err),
		}
	}

	req, err := c.NewRequest(reqBody, "DELETE", "/IAMUser")
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
		keyErr := new(AlksResponseError)
		err = decodeBody(resp, &keyErr)
		if err != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseError, err),
			}
		}

		if keyErr.Errors != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(AlksResponsErrorStrings, strings.Join(keyErr.Errors, ", ")),
			}
		}
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}
	}

	cr := new(DeleteIamUserResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("error parsing DeleteLongTermKeyResponse: %s", err),
		}
	}
	return cr, nil
}

func NewUpdateIamUserRequest(options *IamUserOptions) (*UpdateIamUserRequest, error) {
	if options.IamUserName == nil {
		return nil, fmt.Errorf("IamUserName option must not be nil")
	} else if *options.IamUserName == "" {
		return nil, fmt.Errorf("IamUserName must contain a value")
	}

	iamUser := &UpdateIamUserRequest{}

	if options.Tags != nil {
		iamUser.User.Tags = *options.Tags
	} else {
		return nil, fmt.Errorf("Tags must not be nil on update request, include empty list to remove all non-protected tags")
	}

	return iamUser, nil
}

func (c *Client) UpdateIamUser(options *IamUserOptions) (*UpdateIamUserResponse, *AlksError) {
	request, err := NewUpdateIamUserRequest(options)

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}

	log.Printf("[INFO] update IamUser %s with Tags: %v", *options.IamUserName, *options.Tags)

	accountID, err := c.AccountDetails.GetAccountNumber()
	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error reading Account value: %s", err),
		}
	}

	b, err := json.Marshal(struct {
		UpdateIamUserRequest
	}{*request})

	log.Printf("[INFO] Request Body %v:\n", string(b))

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}
	req, err := c.NewRequest(b, "PATCH", "/iam-users/id/"+accountID+"/"+*options.IamUserName)

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
		updateErr := new(AlksResponseError)
		err = decodeBody(resp, &updateErr)

		if err != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseError, err),
			}
		}

		if updateErr.Errors != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(AlksResponsErrorStrings, strings.Join(updateErr.Errors, ", ")),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}

	}

	respObj := &UpdateIamUserResponse{}
	if err = decodeBody(resp, respObj); err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("error parsing update ltk response: %s", err),
		}
	}

	return respObj, nil
}
