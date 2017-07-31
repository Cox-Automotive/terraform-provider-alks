package alks

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
)

// IamRoleRequqest is used to represent a new IAM Role request.
type IamRoleRequest struct {
	RoleName   string `json:"roleName"`
	RoleType   string `json:"roleType"`
	IncDefPols int    `json:"includeDefaultPolicy"`
}

// IamRoleResponse is used to represent a a IAM Role.
type IamRoleResponse struct {
	RoleName      string   `json:"roleName"`
	RoleType      string   `json:"roleType"`
	RoleArn       string   `json:"roleArn"`
	RoleIPArn     string   `json:"instanceProfileArn"`
	RoleAddedToIP bool     `json:"addedRoleToInstanceProfile"`
	Errors        []string `json:"errors"`
	Exists        bool     `json:"roleExists"`
}

// GetRoleRequest is used to represent a request for details about
// a specific role based on the role's name.
type GetRoleRequest struct {
	RoleName string `json:"roleName"`
}

// DeketeRoleRequest is sued to represent a request for deleting an
// existing IAM role based on the role's name.
type DeleteRoleRequest struct {
	RoleName string `json:"roleName"`
}

// DeleteRoleResponse is used to represent the results of a IAM role
// deletion request.
type DeleteRoleResponse struct {
	RoleName string   `json:"roleName"`
	Status   string   `json:"roleArn"`
	Errors   []string `json:"errors"`
}

// CreateIamRole will create a new IAM role on AWS. If no error is returned
// then you will receive a IamRoleResponse object representing the new role.
func (c *Client) CreateIamRole(roleName string, roleType string, includeDefaultPolicies bool) (*IamRoleResponse, error) {
	log.Printf("[INFO] Creating IAM role: %s", roleName)

	var include int = 0
	if includeDefaultPolicies {
		include = 1
	}

	iam := IamRoleRequest{
		roleName,
		roleType,
		include,
	}

	b, err := json.Marshal(struct {
		IamRoleRequest
		AlksAccount
	}{iam, c.Account})

	if err != nil {
		return nil, fmt.Errorf("Error encoding IAM create role JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/createRole/")
	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return nil, err
	}

	cr := new(IamRoleResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		return nil, fmt.Errorf("Error parsing CreateRole response: %s", err)
	}

	if len(cr.Errors) > 0 {
		return nil, fmt.Errorf("Error creating role: %s", strings.Join(cr.Errors[:], ", "))
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
		AlksAccount
	}{rmRole, c.Account})

	if err != nil {
		return fmt.Errorf("Error encoding IAM delete role JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/deleteRole/")
	if err != nil {
		return err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return err
	}

	del := new(DeleteRoleResponse)
	err = decodeBody(resp, &del)

	if err != nil {
		return fmt.Errorf("Error parsing DeleteRole response: %s", err)
	}

	// TODO you get an error if you delete an already deleted role, need to revist for checking fail/success
	if len(del.Errors) > 0 {
		return fmt.Errorf("Error deleting role: %s", strings.Join(del.Errors[:], ", "))
	}

	return nil
}

// GetIamRole will request the details about an existing IAM role on AWS.
// If no error is returned then you will received a IamRoleResponse object
// representing the existing role. If the role does not exist the IamRoleResponse
// object will also be nil.
func (c *Client) GetIamRole(roleName string) (*IamRoleResponse, error) {
	log.Printf("[INFO] Getting IAM role: %s", roleName)
	getRole := GetRoleRequest{roleName}

	b, err := json.Marshal(struct {
		GetRoleRequest
		AlksAccount
	}{getRole, c.Account})

	if err != nil {
		return nil, fmt.Errorf("Error encoding IAM create role JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/getAccountRole/")
	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return nil, err
	}

	cr := new(IamRoleResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		return nil, fmt.Errorf("Error parsing GetRole response: %s", err)
	}

	if len(cr.Errors) > 0 {
		return nil, fmt.Errorf("Error getting role: %s", strings.Join(cr.Errors[:], ", "))
	}

	if !cr.Exists {
		return nil, fmt.Errorf("Role does not exist.")
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
