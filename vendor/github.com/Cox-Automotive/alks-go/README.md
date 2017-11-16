# alks-go #

alks-go is a Go client library for accessing the ALKS API.

**Documentation:** [![GoDoc](https://godoc.org/github.com/Cox-Automotive/akls-go/github?status.svg)](https://godoc.org/github.com/Cox-Automotive/alks-go)

**Build Status:** [![Build Status](https://travis-ci.org/Cox-Automotive/alks-go.svg?branch=master)](https://travis-ci.org/Cox-Automotive/alks-go)  

alks-go requires Go version 1.7 or greater.

## Usage ##

```go
import "github.com/Cox-Automotive/alks-go"
```

Construct a new ALKS client, then use the various services on the client to
access different parts of the ALKS API. Please note that session creation requires username and password. IAM role CRUD operations can work with either username and password or an STS session. 

*Username/Password Authentication*
```go
client, err := alks.NewClient("http://my.alks.url/rest", "username", "password", "my-acct", "my-role")

// create new STS
resp, err := client.CreateSession(2, false)

log.Printf("Session: %v ~~ %v ~~ %v", resp.AccessKey, resp.SecretKey, resp.SessionToken)
```

*STS Authentication* - Currently only used for IAM role CRUD
```go
client, err := alks.NewSTSClient("http://my.alks.url/rest", "accessKey", "secretKey", "sessionToken", "account")

// create new role
resp, err := client.CreateIamRole("myRole", "Amazon EC2", false)

log.Printf("Role ARN: %v ~~ Role IP ARN: %v", resp.roleArn, resp.roleIPArn)
```

Some API methods don't require an account and role to be provided.
```go
client, err := alks.NewClient("http://my.alks.url/rest", "username", "password", "", "")

// list all available account/roles
resp, err := client.GetAccounts()

for _,acct := range resp.Accounts{
    log.Printf("Account %v Role %v IAM %v", acct.Account, acct.Role, acct.IamActive)
}
```

### Unit Tests ###

You can run the test with Make

```
make test
```