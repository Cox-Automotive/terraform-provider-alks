package alks

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-cleanhttp"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

// AlksAccount is used to represent the configuration for the ALKS client
type AlksAccount struct {
	Username string `json:"userid"`
	Password string `json:"password"`
	Account  string `json:"account"`
	Role     string `json:"role"`
}

type AlksSTS struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Token     string `json:"sessionToken"`
}

// Client represents an ALKS client and contains the account info and base url.
type Client struct {
	Account AlksAccount
	STS     AlksSTS
	BaseURL string

	Http *http.Client
}

// Represents the response from ALKS containing information about a login role
type LoginRoleResponse struct {
	Errors        []string  `json:"errors"`
	StatusMessage string    `json:"statusMessage"`
	RequestId     string    `json:"requestId"`
	LoginRole     LoginRole `json:"loginRole"`
}

// Represents information about a login role
type LoginRole struct {
	Account        string `json:"account"`
	IamKeyActive   bool   `json:"iamKeyActive"`
	MaxKeyDuration int    `json:"maxKeyDuration"`
	Role           string `json:"role"`
}

// NewClient will create a new instance of the ALKS Client. If you don't yet know the account/role
// pass them as nil and then invoke GetAccounts().
func NewClient(url string, username string, password string, account string, role string) (*Client, error) {
	client := Client{
		Account: AlksAccount{
			Username: username,
			Password: password,
			Account:  account,
			Role:     role,
		},
		STS:     AlksSTS{},
		BaseURL: url,
		Http:    cleanhttp.DefaultClient(),
	}

	return &client, nil
}

// NewSTSClient will create a new instance of the ALKS Client using STS tokens.
func NewSTSClient(url string, accessKey string, secretKey string, token string) (*Client, error) {
	client := Client{
		Account: AlksAccount{},
		STS: AlksSTS{
			AccessKey: accessKey,
			SecretKey: secretKey,
			Token:     token,
		},
		BaseURL: url,
		Http:    cleanhttp.DefaultClient(),
	}

	return &client, nil
}

// newRequest will create a new request object for API requests.
func (c *Client) NewRequest(json []byte, method string, endpoint string) (*http.Request, error) {
	u, err := url.Parse(c.BaseURL + endpoint)

	if err != nil {
		return nil, fmt.Errorf("Error parsing base URL: %s", err)
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewBuffer(json))

	if err != nil {
		return nil, fmt.Errorf("Error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "alks-go")

	log.Println("------- ALKS HTTP Request -------")
	requestDump, err := httputil.DumpRequest(req, true)
	if err != nil {
		log.Println(err)
	}
	log.Println(string(requestDump))
	log.Println("-------- !!!!!!!!!! ---------")

	return req, nil
}

// decodeBody will convert a http.Response object to a JSON object.
func decodeBody(resp *http.Response, out interface{}) error {
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	log.Println("------- ALKS HTTP Response -------")
	log.Printf("Status code: %v", resp.StatusCode)
	log.Println(string(body))
	log.Println("-------- !!!!!!!!!! ---------")

	if err = json.Unmarshal(body, &out); err != nil {
		return err
	}

	return nil
}

// checkResp will validate a http.Response based on its status code.
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
	case i == 401: // access denied will still return json
		return resp, nil
	case i == 402:
		return nil, fmt.Errorf("API Error 402: %s", resp.Status)
	case i == 422:
		return nil, fmt.Errorf("API Error 422: %s", resp.Status)
	default:
		return nil, fmt.Errorf("API Error %d: %s", resp.StatusCode, resp.Status)
	}
}

// Durations will provide the valid session durations
func (c *Client) Durations() ([]int, error) {
	log.Printf("[INFO] Requesting allowed durations from ALKS")

	// Use .../me endpoint for getting durations if using STS credentials
	var path string
	if len(strings.TrimSpace(c.Account.Account)) > 0 {
		accountId := c.Account.Account[0:12]
		path = fmt.Sprintf("/loginRoles/id/%v/%v", accountId, c.Account.Role)
	} else {
		path = "/loginRoles/id/me"
	}

	req, err := c.NewRequest(nil, "GET", path)
	if err != nil {
		return nil, err
	}

	resp, err := checkResp(c.Http.Do(req))
	if err != nil {
		return nil, err
	}

	lrr := new(LoginRoleResponse)
	err = decodeBody(resp, &lrr)
	if err != nil {
		return nil, fmt.Errorf("Error parsing LoginRole response: %s", err)
	}

	if len(lrr.Errors) > 0 {
		return nil, fmt.Errorf("Error fetching role information: %s", strings.Join(lrr.Errors[:], ", "))
	}

	maxDuration := lrr.LoginRole.MaxKeyDuration
	durations := make([]int, maxDuration)
	for i := 0; i < maxDuration; i++ {
		durations[i] = i + 1
	}
	return durations, nil
}
