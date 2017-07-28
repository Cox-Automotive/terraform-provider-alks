package alks

import (
	. "github.com/motain/gocheck"
	"time"
)

func (s *S) Test_CreateSession(c *C) {
	testServer.Response(202, nil, sessionCreate)

	resp, err := s.client.CreateSession(2)

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(resp, NotNil)
	c.Assert(resp.AccessKey, Equals, "foo")
	c.Assert(resp.SecretKey, Equals, "bar")
	c.Assert(resp.SessionToken, Equals, "baz")
	c.Assert(resp.SessionDuration, Equals, 2)
	c.Assert(resp.Expires.After(time.Now()), Equals, true)
}

func (s *S) Test_CreateSessionBadTime(c *C) {
	resp, err := s.client.CreateSession(1)

	c.Assert(err, NotNil)
	c.Assert(resp, IsNil)
}

func getIndexByAccount(accounts []AccountRole, account string) (index int) {
	for i, v := range accounts {
		if v.Account == account {
			return i
		}
	}

	return -1
}

func (s *S) Test_GetAccountsPowerUser(c *C) {
	testServer.Response(202, nil, getAccounts)

	resp, err := s.client.GetAccounts()

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(resp, NotNil)
	var index int = getIndexByAccount(resp.Accounts, "123456/ALKSPowerUser - foobarbaz")
	c.Assert(resp.Accounts[index].Account, Equals, "123456/ALKSPowerUser - foobarbaz") // make sure account name is transformed to key
	c.Assert(resp.Accounts[index].Role, Equals, "PowerUser")
	c.Assert(resp.Accounts[index].IamActive, Equals, false)
}

func (s *S) Test_GetAccountsIAMAdmin(c *C) {
	testServer.Response(202, nil, getAccounts)

	resp, err := s.client.GetAccounts()

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(resp, NotNil)
	var index int = getIndexByAccount(resp.Accounts, "234567/ALKSIAMAdmin - foobarbaz2")
	c.Assert(resp.Accounts[index].Account, Equals, "234567/ALKSIAMAdmin - foobarbaz2") // make sure account name is transformed to key
	c.Assert(resp.Accounts[index].Role, Equals, "IAMAdmin")
	c.Assert(resp.Accounts[index].IamActive, Equals, true)
}

func (s *S) Test_GetAccountsAdmin(c *C) {
	testServer.Response(202, nil, getAccounts)

	resp, err := s.client.GetAccounts()

	_ = testServer.WaitRequest()

	c.Assert(err, IsNil)
	c.Assert(resp, NotNil)
	var index int = getIndexByAccount(resp.Accounts, "345678/ALKSAdmin - foobarbaz3")
	c.Assert(resp.Accounts[index].Account, Equals, "345678/ALKSAdmin - foobarbaz3") // make sure account name is transformed to key
	c.Assert(resp.Accounts[index].Role, Equals, "Admin")
	c.Assert(resp.Accounts[index].IamActive, Equals, true)
}

var sessionCreate = `
{
    "accessKey": "foo",
    "secretKey": "bar",
    "sessionToken": "baz"
}
`

// this mapping is so dumb..
var getAccounts = `
{
	"accountListRole": {
		"123456/ALKSPowerUser - foobarbaz": [
		{
			"account": "123456/ALKSPowerUser",
			"role": "PowerUser",
			"iamKeyActive": false
		}
		],
		"234567/ALKSIAMAdmin - foobarbaz2": [
		{
			"account": "234567/ALKSIAMAdmin",
			"role": "IAMAdmin",
			"iamKeyActive": true
		}
		],
		"345678/ALKSAdmin - foobarbaz3": [
		{
			"account": "234567/ALKSAdmin",
			"role": "Admin",
			"iamKeyActive": true
		}
		]
	}
}
`
