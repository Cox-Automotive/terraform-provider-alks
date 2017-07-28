package alks

import (
	"testing"
)

func makeClient(t *testing.T) *Client {
	client, err := NewClient("http://foo.bar.com", "brian", "pass", "acct", "role")

	if err != nil {
		t.Fatalf("err: %v", err)
	}

	if client.BaseURL != "http://foo.bar.com" {
		t.Fatalf("base url not set on client: %s", client.BaseURL)
	}

	if client.Account.Username != "brian" {
		t.Fatalf("account username not set on client: %s", client.Account.Username)
	}

	if client.Account.Password != "pass" {
		t.Fatalf("account password not set on client: %s", client.Account.Password)
	}

	if client.Account.Account != "acct" {
		t.Fatalf("account account not set on client: %s", client.Account.Account)
	}

	if client.Account.Role != "role" {
		t.Fatalf("account role not set on client: %s", client.Account.Role)
	}

	return client
}

func TestClient_NewRequest(t *testing.T) {
	c := makeClient(t)

	json := []byte(`{"fooz":"barz"}`)

	req, err := c.NewRequest(json, "POST", "/endpointfun")
	if err != nil {
		t.Fatalf("bad: %v", err)
	}

	if req.URL.String() != "http://foo.bar.com/endpointfun" {
		t.Fatalf("bad base url: %v", req.URL.String())
	}

	if req.Method != "POST" {
		t.Fatalf("bad method: %v", req.Method)
	}
}
