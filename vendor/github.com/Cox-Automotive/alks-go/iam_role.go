package alks

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// IamRoleRequest is used to represent a new IAM Role request.
type IamRoleRequest struct {
	RoleName       string            `json:"roleName"`
	RoleType       string            `json:"roleType"`
	IncDefPols     int               `json:"includeDefaultPolicy"`
	AlksAccess     bool              `json:"enableAlksAccess"`
	TemplateFields map[string]string `json:"templateFields,omitempty"`
}

// IamTrustRoleRequest is used to represent a new IAM Trust Role request.
type IamTrustRoleRequest struct {
	RoleName   string `json:"roleName"`
	RoleType   string `json:"roleType"`
	TrustArn   string `json:"trustArn"`
	AlksAccess bool   `json:"enableAlksAccess"`
}

// IamRoleResponse is used to represent a a IAM Role.
type IamRoleResponse struct {
	BaseResponse
	RoleName       string            `json:"roleName"`
	RoleType       string            `json:"roleType"`
	RoleArn        string            `json:"roleArn"`
	RoleIPArn      string            `json:"instanceProfileArn"`
	RoleAddedToIP  bool              `json:"addedRoleToInstanceProfile"`
	Exists         bool              `json:"roleExists"`
	TemplateFields map[string]string `json:"templateFields,omitempty"`
}

// GetIamRoleResponse is used to represent a a IAM Role.
type GetIamRoleResponse struct {
	BaseResponse
	RoleName      string `json:"roleName"`
	RoleType      string `json:"roleType"`
	RoleArn       string `json:"roleArn"`
	RoleIPArn     string `json:"instanceProfileArn"`
	RoleAddedToIP bool   `json:"addedRoleToInstanceProfile"`
	Exists        bool   `json:"roleExists"`
	AlksAccess    bool   `json:"machineIdentity"`
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

// CreateIamRole will create a new IAM role on AWS. If no error is returned
// then you will receive a IamRoleResponse object representing the new role.
func (c *Client) CreateIamRole(roleName, roleType string, templateFields map[string]string, includeDefaultPolicies, enableAlksAccess bool) (*IamRoleResponse, error) {
	log.Printf("[INFO] Creating IAM role: %s", roleName)

	var include int
	if includeDefaultPolicies {
		include = 1
	}

	iam := IamRoleRequest{
		roleName,
		roleType,
		include,
		enableAlksAccess,
		templateFields,
	}

	b, err := json.Marshal(struct {
		IamRoleRequest
		AccountDetails
	}{iam, c.AccountDetails})

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
func (c *Client) CreateIamTrustRole(roleName string, roleType string, trustArn string, enableAlksAccess bool) (*IamRoleResponse, error) {
	log.Printf("[INFO] Creating IAM trust role: %s", roleName)

	iam := IamTrustRoleRequest{
		roleName,
		roleType,
		trustArn,
		enableAlksAccess,
	}

	b, err := json.Marshal(struct {
		IamTrustRoleRequest
		AccountDetails
	}{iam, c.AccountDetails})

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

	if !cr.Exists {
		return nil, fmt.Errorf("[%s] Role does not exist", cr.BaseResponse.RequestID)
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
