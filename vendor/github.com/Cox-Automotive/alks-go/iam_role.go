package alks

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// Tag struct is used to represent a AWS Tag
type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AlksError struct {
	StatusMessage string   `json:"statusMessage"`
	Errors        []string `json:"errors"`
	RequestId     string   `json:"requestId"`
}

var ErrorStringFull = "[%s] ALKS Error %d Msg: %s\n Contact the ALKS Team for assistance on Slack at #alks-client-support"
var ErrorStringNoReqId = "ALKS Error %d Msg: %s\n Contact the ALKS Team for assistance on Slack at #alks-client-support"
var ErrorStringOnlyCodeAndReqId = "[%s] ALKS Error %d\n Contact the ALKS Team for assistance on Slack at #alks-client-support"
var ErrorStringOnlyCode = "ALKS Error %d\n Contact the ALKS Team for assistance on Slack at #alks-client-support"
var ParseErrorReqId = "[%s] Error parsing ALKS Error response: %s"
var ParseError = "Error parsing ALKS Error response: %s"

type CreateIamRoleOptions struct {
	RoleName                    *string
	RoleType                    *string
	TrustPolicy                 *map[string]interface{}
	IncludeDefaultPolicies      *bool
	AlksAccess                  *bool
	TrustArn                    *string
	TemplateFields              *map[string]string
	MaxSessionDurationInSeconds *int
	Tags                        *[]Tag
}

// IamRoleRequest is used to represent a new IAM Role request.
type IamRoleRequest struct {
	RoleName                    string                 `json:"roleName"`
	RoleType                    string                 `json:"roleType,omitempty"`
	TrustPolicy                 map[string]interface{} `json:"trustPolicy,omitempty"`
	IncDefPols                  int                    `json:"includeDefaultPolicy,omitempty"`
	AlksAccess                  bool                   `json:"enableAlksAccess,omitempty"`
	TrustArn                    string                 `json:"trustArn,omitempty"`
	TemplateFields              map[string]string      `json:"templateFields,omitempty"`
	MaxSessionDurationInSeconds int                    `json:"maxSessionDurationInSeconds,omitempty"`
	Tags                        []Tag                  `json:"tags,omitempty"`
}

// IamRoleResponse is used to represent a a IAM Role.
type IamRoleResponse struct {
	BaseResponse
	RoleName                    string                 `json:"roleName"`
	RoleType                    string                 `json:"roleType"`
	TrustPolicy                 map[string]interface{} `json:"trustPolicy"`
	RoleArn                     string                 `json:"roleArn"`
	RoleIPArn                   string                 `json:"instanceProfileArn"`
	RoleAddedToIP               bool                   `json:"addedRoleToInstanceProfile"`
	Exists                      bool                   `json:"roleExists"`
	TemplateFields              map[string]string      `json:"templateFields,omitempty"`
	MaxSessionDurationInSeconds int                    `json:"maxSessionDurationInSeconds"`
}

// GetIamRoleResponse is used to represent a a IAM Role.
type GetIamRoleResponse struct {
	BaseResponse
	RoleName                    string                 `json:"roleName"`
	RoleType                    string                 `json:"roleType"`
	TrustPolicy                 map[string]interface{} `json:"trustPolicy"`
	RoleArn                     string                 `json:"roleArn"`
	RoleIPArn                   string                 `json:"instanceProfileArn"`
	RoleAddedToIP               bool                   `json:"addedRoleToInstanceProfile"`
	Exists                      bool                   `json:"roleExists"`
	AlksAccess                  bool                   `json:"machineIdentity"`
	Tags                        []Tag                  `json:"tags"`
	MaxSessionDurationInSeconds int                    `json:"maxSessionDurationInSeconds"`
}

// GetRoleRequest is used to represent a request for details about
// a specific role based on the role's name.
type GetRoleRequest struct {
	RoleName string `json:"roleName"`
}

// DeleteRoleRequest is sued to represent a request for deleting an
// existing IAM role based on the role's name.
type DeleteRoleRequest struct {
	RoleName string `json:"roleName"`
}

// DeleteRoleResponse is used to represent the results of a IAM role
// deletion request.
type DeleteRoleResponse struct {
	BaseResponse
	RoleName string `json:"roleName"`
	Status   string `json:"roleArn"`
}

// AddRoleMachineIdentityRequest is used to represent a request for
// adding a machine identity for a IamRole
type AddRoleMachineIdentityRequest struct {
	RoleARN string `json:"roleARN"`
}

// DeleteRoleMachineIdentityRequest is used to represent a request for
// deleteing a machine identity for a IamRole
type DeleteRoleMachineIdentityRequest struct {
	RoleARN string `json:"roleARN"`
}

// SearchRoleMachineIdentityRequest is used to represent a request for
// searching a machine identity for a given IamRole arn
type SearchRoleMachineIdentityRequest struct {
	RoleARN string `json:"roleARN"`
}

// MachineIdentityResponse is used to represent the results of a add
// machine identity or delete machine identity request.
type MachineIdentityResponse struct {
	BaseResponse
	MachineIdentityArn string `json:"machineIdentityArn"`
}

// Creates a new IamRoleRequest object from options
func NewIamRoleRequest(options *CreateIamRoleOptions) (*IamRoleRequest, error) {
	if options.RoleName == nil {
		return nil, fmt.Errorf("RoleName option must not be nil")
	}

	trustPolicyExists := options.TrustPolicy != nil
	roleTypeExists := options.RoleType != nil
	if trustPolicyExists == roleTypeExists {
		return nil, fmt.Errorf("Either RoleType or TrustPolicy must be included, but not both")
	}

	iam := &IamRoleRequest{
		RoleName: *options.RoleName,
	}

	if roleTypeExists {
		iam.RoleType = *options.RoleType
	}

	if trustPolicyExists {
		iam.TrustPolicy = *options.TrustPolicy
	}

	if options.IncludeDefaultPolicies != nil && *options.IncludeDefaultPolicies {
		iam.IncDefPols = 1
	} else {
		iam.IncDefPols = 0
	}

	if options.AlksAccess != nil {
		iam.AlksAccess = *options.AlksAccess
	} else {
		iam.AlksAccess = false
	}

	if options.TemplateFields != nil {
		iam.TemplateFields = *options.TemplateFields
	} else {
		iam.TemplateFields = nil
	}

	if options.TrustArn != nil {
		iam.TrustArn = *options.TrustArn
	} else {
		iam.TrustArn = ""
	}

	if options.MaxSessionDurationInSeconds != nil {
		iam.MaxSessionDurationInSeconds = *options.MaxSessionDurationInSeconds
	} else {
		iam.MaxSessionDurationInSeconds = 3600
	}

	if options.Tags != nil {
		iam.Tags = *options.Tags
	} else {
		iam.Tags = nil
	}

	return iam, nil
}

// CreateIamRole will create a new IAM role in AWS. If no error is returned
// then you will receive a IamRoleResponse object representing the new role.
func (c *Client) CreateIamRole(options *CreateIamRoleOptions) (*IamRoleResponse, error) {
	request, err := NewIamRoleRequest(options)

	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Creating IAM role: %s", request.RoleName)

	b, err := json.Marshal(struct {
		IamRoleRequest
		AccountDetails
	}{*request, c.AccountDetails})

	if err != nil {
		return nil, fmt.Errorf("Error encoding IAM create role JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/createRole/")
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		createErr := new(AlksError)
		err = decodeBody(resp, &createErr)
		if err != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf(ParseErrorReqId, reqID, err)
			}

			return nil, fmt.Errorf(ParseError, err)
		}

		if createErr.Errors != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf(ErrorStringFull, reqID, resp.StatusCode, createErr.Errors)
			}

			return nil, fmt.Errorf(ErrorStringNoReqId, resp.StatusCode, createErr.Errors)
		}

		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf(ErrorStringOnlyCodeAndReqId, reqID, resp.StatusCode)
		}

		return nil, fmt.Errorf(ErrorStringOnlyCode, resp.StatusCode)
	}

	cr := new(IamRoleResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf("Error parsing CreateRole response: [%s] %s", reqID, err)
		}

		return nil, fmt.Errorf("Error parsing CreateRole response: %s", err)
	}

	if cr.RequestFailed() {
		return nil, fmt.Errorf("Error creating role: [%s] %s", cr.BaseResponse.RequestID, strings.Join(cr.GetErrors(), ", "))
	}

	return cr, nil
}

// CreateIamTrustRole will create a new IAM trust role on AWS. If no error is returned
// then you will receive a IamRoleResponse object representing the new role.
func (c *Client) CreateIamTrustRole(options *CreateIamRoleOptions) (*IamRoleResponse, error) {
	request, err := NewIamRoleRequest(options)

	b, err := json.Marshal(struct {
		IamRoleRequest
		AccountDetails
	}{*request, c.AccountDetails})

	if err != nil {
		return nil, fmt.Errorf("Error encoding IAM create trust role JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/createNonServiceRole/")
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		trustErr := new(AlksError)
		err = decodeBody(resp, &trustErr)

		if err != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf(ParseErrorReqId, reqID, err)
			}

			return nil, fmt.Errorf(ParseError, err)
		}

		if trustErr.Errors != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf(ErrorStringFull, reqID, resp.StatusCode, trustErr.Errors)
			}

			return nil, fmt.Errorf(ErrorStringNoReqId, resp.StatusCode, trustErr.Errors)
		}

		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf(ErrorStringOnlyCodeAndReqId, reqID, resp.StatusCode)
		}

		return nil, fmt.Errorf(ErrorStringOnlyCode, resp.StatusCode)
	}

	cr := new(IamRoleResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf("Error parsing CreateTrustRole response: [%s] %s", reqID, err)
		}

		return nil, fmt.Errorf("Error parsing CreateTrustRole response: %s", err)
	}

	if cr.RequestFailed() {
		return nil, fmt.Errorf("Error creating trust role: [%s] %s", cr.BaseResponse.RequestID, strings.Join(cr.GetErrors(), ", "))
	}

	return cr, nil
}

type UpdateIamRoleRequest struct {
	RoleName *string `json:"roleName"`
	Tags     *[]Tag  `json:"tags"`
}

type UpdateIamRoleResponse struct {
	BaseResponse
	RoleArn         *string `json:"roleArn"`
	RoleName        *string `json:"roleName"`
	BasicAuth       *bool   `json:"basicAuthUsed"`
	Exists          *bool   `json:"roleExists"`
	RoleIPArn       *string `json:"instanceProfileArn"`
	MachineIdentity *bool   `json:"isMachineIdentity"`
	Tags            *[]Tag  `json:"tags"`
}

/* UpdateIamRole adds resource tags to an existing IAM role.
 */
func (c *Client) UpdateIamRole(options *UpdateIamRoleRequest) (*UpdateIamRoleResponse, error) {
	if err := options.updateIamRoleValidate(); err != nil {
		return nil, err
	}
	log.Printf("[INFO] update IAM role %s with Tags: %v", *options.RoleName, *options.Tags)

	b, err := json.Marshal(struct {
		UpdateIamRoleRequest
		AccountDetails
	}{*options, c.AccountDetails})
	if err != nil {
		return nil, err
	}
	req, err := c.NewRequest(b, "PATCH", "/role/")
	if err != nil {
		return nil, err
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		updateErr := new(AlksError)
		err = decodeBody(resp, &updateErr)

		if err != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf(ParseErrorReqId, reqID, err)
			}

			return nil, fmt.Errorf(ParseError, err)
		}

		if updateErr.Errors != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf(ErrorStringFull, reqID, resp.StatusCode, updateErr.Errors)
			}

			return nil, fmt.Errorf(ErrorStringNoReqId, resp.StatusCode, updateErr.Errors)
		}

		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf(ErrorStringOnlyCodeAndReqId, reqID, resp.StatusCode)
		}

		return nil, fmt.Errorf(ErrorStringOnlyCode, resp.StatusCode)
	}

	respObj := &UpdateIamRoleResponse{}
	if err = decodeBody(resp, respObj); err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf("error parsing update role response: [%s] %s", reqID, err)
		}
		return nil, fmt.Errorf("error parsing update role response: %s", err)
	}
	if respObj.RequestFailed() {
		return nil, fmt.Errorf("error from update IAM role request: [%s] %s", respObj.RequestID, strings.Join(respObj.GetErrors(), ", "))
	}

	return respObj, nil
}

func (req *UpdateIamRoleRequest) updateIamRoleValidate() error {
	if req.RoleName == nil {
		return fmt.Errorf("roleName option must not be nil")
	}
	if req.Tags == nil {
		return fmt.Errorf("tags option must not be nil")
	}
	return nil
}

// DeleteIamRole will delete an existing IAM role from AWS. If no error is returned
// then the deletion was successful.
func (c *Client) DeleteIamRole(id string) error {
	log.Printf("[INFO] Deleting IAM role: %s", id)

	rmRole := DeleteRoleRequest{id}

	b, err := json.Marshal(struct {
		DeleteRoleRequest
		AccountDetails
	}{rmRole, c.AccountDetails})

	if err != nil {
		return fmt.Errorf("Error encoding IAM delete role JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/deleteRole/")
	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		delErr := new(AlksError)
		err = decodeBody(resp, &delErr)

		if err != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return fmt.Errorf(ParseErrorReqId, reqID, err)
			}

			return fmt.Errorf(ParseError, err)
		}

		if delErr.Errors != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return fmt.Errorf(ErrorStringFull, reqID, resp.StatusCode, delErr.Errors)
			}

			return fmt.Errorf(ErrorStringNoReqId, resp.StatusCode, delErr.Errors)
		}

		if reqID := GetRequestID(resp); reqID != "" {
			return fmt.Errorf(ErrorStringOnlyCodeAndReqId, reqID, resp.StatusCode)
		}

		return fmt.Errorf(ErrorStringOnlyCode, resp.StatusCode)
	}

	del := new(DeleteRoleResponse)
	err = decodeBody(resp, &del)

	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return fmt.Errorf("Error parsing DeleteRole response: [%s] %s", reqID, err)
		}

		return fmt.Errorf("Error parsing DeleteRole response: %s", err)
	}

	// TODO you get an error if you delete an already deleted role, need to revist for checking fail/success
	if del.RequestFailed() {
		return fmt.Errorf("Error deleting role: [%s] %s", del.BaseResponse.RequestID, strings.Join(del.GetErrors(), ", "))
	}

	return nil
}

// GetIamRole will request the details about an existing IAM role on AWS.
// If no error is returned then you will received a IamRoleResponse object
// representing the existing role. If the role does not exist the IamRoleResponse
// object will also be nil.
func (c *Client) GetIamRole(roleName string) (*GetIamRoleResponse, error) {
	log.Printf("[INFO] Getting IAM role: %s", roleName)
	getRole := GetRoleRequest{roleName}

	b, err := json.Marshal(struct {
		GetRoleRequest
		AccountDetails
	}{getRole, c.AccountDetails})

	if err != nil {
		return nil, fmt.Errorf("Error encoding IAM create role JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/getAccountRole/")
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if (resp.StatusCode < 200 || resp.StatusCode >= 300) && resp.StatusCode != 404 {
		getErr := new(AlksError)
		err = decodeBody(resp, &getErr)

		if err != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf(ParseErrorReqId, reqID, err)
			}

			return nil, fmt.Errorf(ParseError, err)
		}

		if getErr.Errors != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf(ErrorStringFull, reqID, resp.StatusCode, getErr.Errors)
			}

			return nil, fmt.Errorf(ErrorStringNoReqId, resp.StatusCode, getErr.Errors)
		}

		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf(ErrorStringOnlyCodeAndReqId, reqID, resp.StatusCode)
		}

		return nil, fmt.Errorf(ErrorStringOnlyCode, resp.StatusCode)
	}

	cr := new(GetIamRoleResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf("Error parsing GetRole response: [%s] %s", reqID, err)
		}

		return nil, fmt.Errorf("Error parsing GetRole response: %s", err)
	}

	if cr.RequestFailed() {
		return nil, fmt.Errorf("Error getting role: [%s] %s", cr.BaseResponse.RequestID, strings.Join(cr.GetErrors(), ", "))
	}

	// This is here because ALKS returns a string representation of a Java array
	// with the only entry being the instance profile ARN (ie: "[\"ARN\"]")
	// A simple regex fixes the formatting issue and using existing instance
	// profiles works again. Every IAM role doesn't return an instance profile,
	// so we have to make sure the string isn't empty.
	if len(cr.RoleIPArn) > 0 {
		re := regexp.MustCompile("^\\[\\\"(.+)\\\"\\]$")
		matches := re.FindStringSubmatch(cr.RoleIPArn)
		if len(matches) > 1 {
			cr.RoleIPArn = matches[1]
		}
	}

	return cr, nil
}

// AddRoleMachineIdentity enable machine identity for a IamRole.
// If no error is returned then you will receieve the arn for the machine identity that was created.
func (c *Client) AddRoleMachineIdentity(roleARN string) (*MachineIdentityResponse, error) {
	log.Printf("[INFO] Adding role machine identity: %s", roleARN)
	addMI := AddRoleMachineIdentityRequest{roleARN}

	b, err := json.Marshal(struct {
		AddRoleMachineIdentityRequest
	}{addMI})

	if err != nil {
		return nil, fmt.Errorf("Error encoding add role machine identity JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/roleMachineIdentity/")
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		addErr := new(AlksError)
		err = decodeBody(resp, &addErr)

		if err != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf(ParseErrorReqId, reqID, err)
			}

			return nil, fmt.Errorf(ParseError, err)
		}

		if addErr.Errors != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf(ErrorStringFull, reqID, resp.StatusCode, addErr.Errors)
			}

			return nil, fmt.Errorf(ErrorStringNoReqId, resp.StatusCode, addErr.Errors)
		}

		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf(ErrorStringOnlyCodeAndReqId, reqID, resp.StatusCode)
		}

		return nil, fmt.Errorf(ErrorStringOnlyCode, resp.StatusCode)
	}

	cr := new(MachineIdentityResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf("Error parsing MachineIdentitiyResponse: [%s] %s", reqID, err)
		}

		return nil, fmt.Errorf("Error parsing MachineIdentityResponse: %s", err)
	}

	if cr.RequestFailed() {
		return nil, fmt.Errorf("Error creating machine identity: [%s] %s", cr.BaseResponse.RequestID, strings.Join(cr.GetErrors(), ", "))
	}

	return cr, nil
}

// DeleteRoleMachineIdentity disable machine identity for a IamRole.
// If no error is returned then you will receieve the arn for the machine identity that was deleted.
func (c *Client) DeleteRoleMachineIdentity(roleARN string) (*MachineIdentityResponse, error) {
	log.Printf("[INFO] Deleting role machine identity: %s", roleARN)
	deleteMI := DeleteRoleMachineIdentityRequest{roleARN}

	b, err := json.Marshal(struct {
		DeleteRoleMachineIdentityRequest
	}{deleteMI})

	if err != nil {
		return nil, fmt.Errorf("Error encoding delete role machine identity JSON: %s", err)
	}

	req, err := c.NewRequest(b, "DELETE", "/roleMachineIdentity/")
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		delErr := new(AlksError)
		err = decodeBody(resp, &delErr)

		if err != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf(ParseErrorReqId, reqID, err)
			}

			return nil, fmt.Errorf(ParseError, err)
		}

		if delErr.Errors != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf(ErrorStringFull, reqID, resp.StatusCode, delErr.Errors)
			}

			return nil, fmt.Errorf(ErrorStringNoReqId, resp.StatusCode, delErr.Errors)
		}

		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf(ErrorStringOnlyCodeAndReqId, reqID, resp.StatusCode)
		}

		return nil, fmt.Errorf(ErrorStringOnlyCode, resp.StatusCode)
	}

	dr := new(MachineIdentityResponse)
	err = decodeBody(resp, &dr)

	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf("Error parsing MachineIdentityResponse: [%s] %s", reqID, err)
		}

		return nil, fmt.Errorf("Error parsing MachineIdenttiyResponse: %s", err)
	}

	if dr.RequestFailed() {
		return nil, fmt.Errorf("Error deleting machine identity: [%s] %s", dr.BaseResponse.RequestID, strings.Join(dr.GetErrors(), ", "))
	}

	return dr, nil
}

// SearchRoleMachineIdentity searches for a machine identity for a given roleARN
// If no error is returned then you will receive the arn of the machine identity for the given roleARN
func (c *Client) SearchRoleMachineIdentity(roleARN string) (*MachineIdentityResponse, error) {
	log.Printf("[INFO] Searching role machine identity: %s", roleARN)
	searchMI := SearchRoleMachineIdentityRequest{roleARN}

	b, err := json.Marshal(struct {
		SearchRoleMachineIdentityRequest
	}{searchMI})

	if err != nil {
		return nil, fmt.Errorf("Error decoding search role machine identity JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/roleMachineIdentity/search/")
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		searchErr := new(AlksError)
		err = decodeBody(resp, &searchErr)

		if err != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf(ParseErrorReqId, reqID, err)
			}

			return nil, fmt.Errorf(ParseError, err)
		}

		if searchErr.Errors != nil {
			if reqID := GetRequestID(resp); reqID != "" {
				return nil, fmt.Errorf("[%s] ALKS Error Code: %d Msg: %s\n Contact the ALKS Team for assistance on Slack at #alks-client-support", reqID, resp.StatusCode, searchErr.Errors)
			}

			return nil, fmt.Errorf("ALKS Error Code: %d Msg: %s\n Contact the ALKS Team for assistance on Slack at #alks-client-support", resp.StatusCode, searchErr.Errors)
		}

		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf(ErrorStringOnlyCodeAndReqId, reqID, resp.StatusCode)
		}

		return nil, fmt.Errorf(ErrorStringOnlyCode, resp.StatusCode)
	}

	sr := new(MachineIdentityResponse)
	err = decodeBody(resp, &sr)

	if err != nil {
		if reqID := GetRequestID(resp); reqID != "" {
			return nil, fmt.Errorf("Error parsing MachineIdentityResponse: [%s] %s", reqID, err)
		}

		return nil, fmt.Errorf("Error parsing MachineIdentityResponse: %s", err)
	}

	if sr.RequestFailed() {
		return nil, fmt.Errorf("Error searching machine identity [%s] %s", sr.BaseResponse.RequestID, strings.Join(sr.GetErrors(), ", "))
	}

	return sr, nil
}
