package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform/helper/schema"
)

type AlksClient struct {
	BaseURL  string
	Username string
	Password string
	Account  string
	Role     string

	Http *http.Client
}

type CreateIamKeyReq struct {
	Username    string `json:"userid"`
	Password    string `json:"password"`
	Account     string `json:"account"`
	Role        string `json:"role"`
	SessionTime int    `json:"sessionTime"`
}

type CreateIamRoleReq struct {
	Username   string `json:"userid"`
	Password   string `json:"password"`
	Account    string `json:"account"`
	Role       string `json:"role"`
	RoleName   string `json:"roleName"`
	RoleType   string `json:"roleType"`
	IncDefPols int    `json:"includeDefaultPolicy"`
}

type StsResponse struct {
	AccessKey    string `json:"accessKey"`
	SessionKey   string `json:"secretKey"`
	SessionToken string `json:"sessionToken"`
}

type CreateRoleResponse struct {
	RoleName      string   `json:"roleName"`
	RoleType      string   `json:"roleType"`
	RoleArn       string   `json:"roleArn"`
	RoleIPArn     string   `json:"instanceProfileArn"`
	RoleAddedToIP bool     `json:"addedRoleToInstanceProfile"`
	Errors        []string `json:"errors"`
}

type GetRoleResponse struct {
	RoleName   string `json:"roleName"`
	RoleType   string `json:"roleType"`
	RoleArn    string `json:"roleArn"`
	RoleIPArn  string `json:"instanceProfileArn"`
	IncDefPols int    `json:"includeDefaultPolicy"`
}

func NewAlksClient(url string, username string, password string, account string, role string) (*AlksClient, error) {
	alksClient := AlksClient{
		BaseURL:  url,
		Username: username,
		Password: password,
		Account:  account,
		Role:     role,
		Http:     cleanhttp.DefaultClient(),
	}

	return &alksClient, nil
}

func (c *AlksClient) NewRequest(json []byte, method string, endpoint string) (*http.Request, error) {
	u, err := url.Parse(c.BaseURL + endpoint)

	if err != nil {
		return nil, fmt.Errorf("Error parsing base URL: %s", err)
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(json))

	if err != nil {
		return nil, fmt.Errorf("Error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	if err = json.Unmarshal(body, &out); err != nil {
		return err
	}

	return nil
}

func checkResp(resp *http.Response, err error) (*http.Response, error) {
	if err != nil {
		return resp, err
	}

	switch i := resp.StatusCode; {
	case i == 200:
		return resp, nil
	case i == 201:
		return resp, nil
	case i == 202:
		return resp, nil
	case i == 204:
		return resp, nil
	case i == 400:
		return nil, fmt.Errorf("API Error 400: %s", resp.Status)
	case i == 401:
		return nil, fmt.Errorf("API Error 401: %s", resp.Status)
	case i == 402:
		return nil, fmt.Errorf("API Error 402: %s", resp.Status)
	case i == 422:
		return nil, fmt.Errorf("API Error 422: %s", resp.Status)
	default:
		return nil, fmt.Errorf("API Error: %s", resp.Status)
	}
}

func (c *AlksClient) CreateIamKey() (*StsResponse, error) {

	iam := CreateIamKeyReq{c.Username, c.Password, c.Account, c.Role, 1}
	b, err := json.Marshal(iam)

	if err != nil {
		return nil, fmt.Errorf("Error encoding IAM create key JSON: %s", err)
	}

	req, err := c.NewRequest(b, "POST", "/getIAMKeys/")
	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return nil, err
	}

	sts := new(StsResponse)
	err = decodeBody(resp, &sts)

	if err != nil {
		return nil, fmt.Errorf("Error parsing STS response: %s", err)
	}

	return sts, nil
}

func (c *AlksClient) CreateIamRole(roleName string, roleType string, includeDefaultPolicies bool) (*CreateRoleResponse, error) {
	var include int = 0
	if includeDefaultPolicies {
		include = 1
	}

	iam := CreateIamRoleReq{c.Username, c.Password, c.Account, c.Role, roleName, roleType, include}
	b, err := json.Marshal(iam)

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

	cr := new(CreateRoleResponse)
	err = decodeBody(resp, &cr)

	if err != nil {
		return nil, fmt.Errorf("Error parsing CreateRole response: %s", err)
	}

	if len(cr.Errors) > 0 {
		return nil, fmt.Errorf("Error creating role: %s", strings.Join(cr.Errors[:], ", "))
	}

	return cr, nil
}

func DeleteIamRole(d *schema.ResourceData, meta interface{}) error {
	// client := meta.(*AlksClient)

	log.Printf("[INFO] Deleting IAM role: %s", d.Id())

	// TODO: IMPLEMENT

	return nil
}

func (c *AlksClient) GetIamRole(arn string) (*GetRoleResponse, error) {
	//  TODO: IMPLEMENT
	return nil, nil
}

func (c *AlksClient) GetIamRoleByName(roleName string) (*GetRoleResponse, error) {
	//  TODO: IMPLEMENT
	return nil, nil
}
