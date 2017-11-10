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
	Account   string `json:"account"`
	Role      string `json:"role"`
}

// Client represents an ALKS client and contains the account info and base url.
type Client struct {
	Account AlksAccount
	STS     AlksSTS
	BaseURL string

	Http *http.Client
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
func NewSTSClient(url string, accessKey string, secretKey string, token string, account string) (*Client, error) {
	client := Client{
		Account: AlksAccount{},
		STS: AlksSTS{
			AccessKey: accessKey,
			SecretKey: secretKey,
			Token:     token,
			Account:   account,
			Role:      "Admin", // TEMPORARY: FIXME: TODO: Should be removed once the API req goes away
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
		return nil, fmt.Errorf("API Error: Please validate your username/password and account/role.")
	}
}

// Durations will provide the valid session durations
func (c *Client) Durations() []int {
	return []int{1, 2, 6, 12, 18}
}
