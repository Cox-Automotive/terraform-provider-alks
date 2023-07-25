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
func (c *Client) CreateIamRole(options *CreateIamRoleOptions) (*IamRoleResponse, *AlksError) {
	request, err := NewIamRoleRequest(options)

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}

	log.Printf("[INFO] Creating IAM role: %s", request.RoleName)

	b, err := json.Marshal(struct {
		IamRoleRequest
		AccountDetails
	}{*request, c.AccountDetails})

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error encoding IAM create role JSON: %s", err),
		}
	}

	req, err := c.NewRequest(b, "POST", "/createRole/")
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
		alksResponseErr := new(AlksResponseError)
		err = decodeBody(resp, &alksResponseErr)

		if err != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseErrorReqId, reqID, err),
			}
		}

		if alksResponseErr.Errors != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(AlksResponsErrorStrings, strings.Join(alksResponseErr.Errors, ", ")),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}
	}

	cr := new(IamRoleResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error parsing CreateRole response: [%s] %s", reqID, err),
		}
	}

	if cr.RequestFailed() {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error creating role: [%s] %s", cr.BaseResponse.RequestID, strings.Join(cr.GetErrors(), ", ")),
		}
	}

	return cr, nil
}

// CreateIamTrustRole will create a new IAM trust role on AWS. If no error is returned
// then you will receive a IamRoleResponse object representing the new role.
func (c *Client) CreateIamTrustRole(options *CreateIamRoleOptions) (*IamRoleResponse, *AlksError) {
	request, err := NewIamRoleRequest(options)

	b, err := json.Marshal(struct {
		IamRoleRequest
		AccountDetails
	}{*request, c.AccountDetails})

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error encoding IAM create trust role JSON: %s", err),
		}
	}

	req, err := c.NewRequest(b, "POST", "/createNonServiceRole/")
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
		trustErr := new(AlksResponseError)
		err = decodeBody(resp, &trustErr)

		if err != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseErrorReqId, reqID, err),
			}
		}

		if trustErr.Errors != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(AlksResponsErrorStrings, strings.Join(trustErr.Errors, ", ")),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}
	}

	cr := new(IamRoleResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error parsing CreateTrustRole response: [%s] %s", reqID, err),
		}
	}

	if cr.RequestFailed() {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error creating trust role: [%s] %s", cr.BaseResponse.RequestID, strings.Join(cr.GetErrors(), ", ")),
		}
	}

	return cr, nil
}

type UpdateIamRoleRequest struct {
	RoleName    *string                 `json:"roleName"`
	Tags        *[]Tag                  `json:"tags"`
	TrustPolicy *map[string]interface{} `json:"trustPolicy"`
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

// Updates an IAM role with the given options.
func (c *Client) UpdateIamRole(options *UpdateIamRoleRequest) (*UpdateIamRoleResponse, *AlksError) {
	if err := options.updateIamRoleValidate(); err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}
	// considering a non empty tag object
	if options.Tags != nil {
		log.Printf("[INFO] update IAM role %s with tags: %v", *options.RoleName, *options.Tags)
	}
	// considering a non empty TrustPolicy map
	if options.TrustPolicy != nil {
		log.Printf("[INFO] update IAM role %s with trust policy: %v", *options.RoleName, *options.TrustPolicy)
	}

	b, err := json.Marshal(struct {
		UpdateIamRoleRequest
		AccountDetails
	}{*options, c.AccountDetails})
	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}
	req, err := c.NewRequest(b, "PATCH", "/role/")
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

	respObj := &UpdateIamRoleResponse{}
	if err = decodeBody(resp, respObj); err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error parsing updateRole response"),
		}
	}
	if respObj.RequestFailed() {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  respObj.RequestID,
			Err:        fmt.Errorf("Error from update IAM role request: %s", strings.Join(respObj.GetErrors(), ", ")),
		}
	}
	return respObj, nil
}

func (req *UpdateIamRoleRequest) updateIamRoleValidate() error {
	if req.RoleName == nil {
		return fmt.Errorf("roleName option must not be nil")
	}
	return nil
}

// DeleteIamRole will delete an existing IAM role from AWS. If no error is returned
// then the deletion was successful.
func (c *Client) DeleteIamRole(id string) *AlksError {
	log.Printf("[INFO] Deleting IAM role: %s", id)

	rmRole := DeleteRoleRequest{id}

	b, err := json.Marshal(struct {
		DeleteRoleRequest
		AccountDetails
	}{rmRole, c.AccountDetails})

	if err != nil {
		return &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error encoding IAM delete role JSON: %s", err),
		}
	}

	req, err := c.NewRequest(b, "POST", "/deleteRole/")
	if err != nil {
		return &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        err,
		}
	}

	reqID := GetRequestID(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		delErr := new(AlksResponseError)
		err = decodeBody(resp, &delErr)

		if err != nil {
			return &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseError, err),
			}
		}

		if delErr.Errors != nil {
			return &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(AlksResponsErrorStrings, strings.Join(delErr.Errors, ", ")),
			}
		}

		return &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}
	}

	del := new(DeleteRoleResponse)
	err = decodeBody(resp, &del)

	if err != nil {
		return &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error parsing deleteRole response: %s", err),
		}
	}

	// TODO you get an error if you delete an already deleted role, need to revist for checking fail/success
	if del.RequestFailed() {
		return &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  del.BaseResponse.RequestID,
			Err:        fmt.Errorf("Error deleting role: %s", strings.Join(del.GetErrors(), ", ")),
		}
	}

	return nil
}

// GetIamRole will request the details about an existing IAM role on AWS.
// If no error is returned then you will received a IamRoleResponse object
// representing the existing role. If the role does not exist the IamRoleResponse
// object will also be nil.
func (c *Client) GetIamRole(roleName string) (*GetIamRoleResponse, *AlksError) {
	log.Printf("[INFO] Getting IAM role: %s", roleName)
	getRole := GetRoleRequest{roleName}

	b, err := json.Marshal(struct {
		GetRoleRequest
		AccountDetails
	}{getRole, c.AccountDetails})

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error encoding IAM get role JSON: %s", err),
		}
	}

	req, err := c.NewRequest(b, "POST", "/getAccountRole/")
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
		getErr := new(AlksResponseError)
		err = decodeBody(resp, &getErr)

		if err != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseError, err),
			}
		}

		if getErr.Errors != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(AlksResponsErrorStrings, strings.Join(getErr.Errors, ", ")),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}
	}

	cr := new(GetIamRoleResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error parsing getRole response: %s", err),
		}
	}

	if cr.RequestFailed() {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  cr.BaseResponse.RequestID,
			Err:        fmt.Errorf("Error getting role: %s", strings.Join(cr.GetErrors(), ", ")),
		}
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
func (c *Client) AddRoleMachineIdentity(roleARN string) (*MachineIdentityResponse, *AlksError) {
	log.Printf("[INFO] Adding role machine identity: %s", roleARN)
	addMI := AddRoleMachineIdentityRequest{roleARN}

	b, err := json.Marshal(struct {
		AddRoleMachineIdentityRequest
	}{addMI})

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error encoding add role machine identity JSON: %s", err),
		}
	}

	req, err := c.NewRequest(b, "POST", "/roleMachineIdentity/")
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
		addErr := new(AlksResponseError)
		err = decodeBody(resp, &addErr)

		if err != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseError, err),
			}
		}

		if addErr.Errors != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(AlksResponsErrorStrings, strings.Join(addErr.Errors, ", ")),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}
	}

	cr := new(MachineIdentityResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error parsing MachineIdentitiyResponse response: %s", err),
		}
	}

	if cr.RequestFailed() {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error creating machine identity: %s", strings.Join(cr.GetErrors(), ", ")),
		}
	}

	return cr, nil
}

// DeleteRoleMachineIdentity disable machine identity for a IamRole.
// If no error is returned then you will receieve the arn for the machine identity that was deleted.
func (c *Client) DeleteRoleMachineIdentity(roleARN string) (*MachineIdentityResponse, *AlksError) {
	log.Printf("[INFO] Deleting role machine identity: %s", roleARN)
	deleteMI := DeleteRoleMachineIdentityRequest{roleARN}

	b, err := json.Marshal(struct {
		DeleteRoleMachineIdentityRequest
	}{deleteMI})

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error encoding delete role machine identity JSON: %s", err),
		}
	}

	req, err := c.NewRequest(b, "DELETE", "/roleMachineIdentity/")
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
		delErr := new(AlksResponseError)
		err = decodeBody(resp, &delErr)

		if err != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseError, err),
			}
		}

		if delErr.Errors != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(AlksResponsErrorStrings, strings.Join(delErr.Errors, ", ")),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}
	}

	dr := new(MachineIdentityResponse)
	err = decodeBody(resp, &dr)

	if err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error parsing machineIdentity response: %s", err),
		}
	}

	if dr.RequestFailed() {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error deleting machine identity: %s", strings.Join(dr.GetErrors(), ", ")),
		}
	}

	return dr, nil
}

// SearchRoleMachineIdentity searches for a machine identity for a given roleARN
// If no error is returned then you will receive the arn of the machine identity for the given roleARN
func (c *Client) SearchRoleMachineIdentity(roleARN string) (*MachineIdentityResponse, *AlksError) {
	log.Printf("[INFO] Searching role machine identity: %s", roleARN)
	searchMI := SearchRoleMachineIdentityRequest{roleARN}

	b, err := json.Marshal(struct {
		SearchRoleMachineIdentityRequest
	}{searchMI})

	if err != nil {
		return nil, &AlksError{
			StatusCode: 0,
			RequestId:  "",
			Err:        fmt.Errorf("Error decoding search role machine identity JSON: %s", err),
		}
	}

	req, err := c.NewRequest(b, "POST", "/roleMachineIdentity/search/")
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
		searchErr := new(AlksResponseError)
		err = decodeBody(resp, &searchErr)

		if err != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(ParseError, err),
			}
		}

		if searchErr.Errors != nil {
			return nil, &AlksError{
				StatusCode: resp.StatusCode,
				RequestId:  reqID,
				Err:        fmt.Errorf(AlksResponsErrorStrings, strings.Join(searchErr.Errors, ", ")),
			}
		}

		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf(GenericAlksError),
		}
	}

	sr := new(MachineIdentityResponse)
	err = decodeBody(resp, &sr)

	if err != nil {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  reqID,
			Err:        fmt.Errorf("Error parsing MachineIdentity response: %s", err),
		}
	}

	if sr.RequestFailed() {
		return nil, &AlksError{
			StatusCode: resp.StatusCode,
			RequestId:  sr.BaseResponse.RequestID,
			Err:        fmt.Errorf("Error searching machine identity %s", strings.Join(sr.GetErrors(), ", ")),
		}
	}

	return sr, nil
}
