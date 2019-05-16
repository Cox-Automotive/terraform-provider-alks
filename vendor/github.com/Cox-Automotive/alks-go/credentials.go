package alks

import (
	"errors"
	"net/http"
)

const (
	accessKeyHeader    = "ALKS-STS-Access-Key"
	secretKeyHeader    = "ALKS-STS-Secret-Key"
	sessionTokenHeader = "ALKS-STS-Session-Token"
)

// Basic represents LDAP based credentials in the configuration of the ALKS client
type Basic struct {
	Username string `json:"-"`
	Password string `json:"-"`
}

// STS represents AWS STS credentials in the configuration of the ALKS client
type STS struct {
	AccessKey    string `json:"-"`
	SecretKey    string `json:"-"`
	SessionToken string `json:"-"`
}

// Bearer represents an Okta bearer token in the configuration of the ALKS client
type Bearer struct {
	Token string `json:"-"`
}

// AuthInjecter is the interface that wraps the InjectAuth method.
//
// Implementations are expect to add their authentication data to request without
// destroying existing data (if any) and should implement fallbacks when
// possible.  Failing that, an error should be reported to the caller.
type AuthInjecter interface {
	InjectAuth(req *http.Request) error
}

// InjectAuth will add an Authorization header to an ALKS client request containing
// the caller's username and password.
func (b *Basic) InjectAuth(req *http.Request) error {
	if _, _, ok := req.BasicAuth(); ok {
		return errors.New("Basic Auth header already exists")
	}

	req.SetBasicAuth(b.Username, b.Password)

	return nil
}

// InjectAuth will add ALKS headers to client requests containing
// the caller's STS credentials.
func (s *STS) InjectAuth(req *http.Request) error {
	if req.Header.Get(accessKeyHeader) != "" &&
		req.Header.Get(secretKeyHeader) != "" &&
		req.Header.Get(sessionTokenHeader) != "" {
		return errors.New("STS Auth headers already exist")
	}

	req.Header.Add(accessKeyHeader, s.AccessKey)
	req.Header.Add(secretKeyHeader, s.SecretKey)
	req.Header.Add(sessionTokenHeader, s.SessionToken)

	return nil
}

// InjectAuth will add an authorization header to an ALKS client request containing
// the caller's Okta bearer token.
func (b *Bearer) InjectAuth(req *http.Request) error {
	if req.Header.Get("Authorization") != "" {
		return errors.New("Authorization header already exists")
	}

	req.Header.Add("Authorization", "Bearer "+b.Token)

	return nil
}
