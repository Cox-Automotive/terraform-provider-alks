package alks

import (
	"github.com/Cox-Automotive/alks-go/testutils"
	"testing"

	. "github.com/motain/gocheck"
)

func Test(t *testing.T) {
	TestingT(t)
}

type S struct {
	client *Client
}

var _ = Suite(&S{})

var testServer = testutil.NewHTTPServer()

func (s *S) SetUpSuite(c *C) {
	testServer.Start()
	var err error
	s.client, err = NewClient("http://localhost:4200", "brian", "pass", "acct", "role")
	if err != nil {
		panic(err)
	}
}

func (s *S) TearDownTest(c *C) {
	testServer.Flush()
}

func (s *S) Test_CreateIamSession(c *C) {
	testServer.Response(202, nil, iamResponse)

	resp, err := s.client.CreateIamSession()

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(resp.AccessKey, Equals, "thisismykey")
	c.Assert(resp.SessionKey, Equals, "thisismysecret")
	c.Assert(resp.SessionToken, Equals, "thisismysession")
}

var iamResponse = `
{
    "accessKey": "thisismykey",
    "secretKey": "thisismysecret",
    "sessionToken": "thisismysession"
}
`
