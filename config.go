package main

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"

	"github.com/Cox-Automotive/alks-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Version number, to be injected at link time
// to set, add `-ldflags "-X main.versionNumber=1.2.3"` to the go build command
var versionNumber string
var ErrNoValidCredentialSources = errors.New(`No valid credential sources found for ALKS Provider.
Please see https://github.com/Cox-Automotive/terraform-provider-alks#authentication for more information on
providing credentials for the ALKS Provider`)

// Config stores ALKS configuration and credentials
type Config struct {
	URL           string
	AccessKey     string
	SecretKey     string
	Token         string
	CredsFilename string
	Profile       string
	AssumeRole    assumeRoleDetails
}

type assumeRoleDetails struct {
	RoleARN     string
	SessionName string
	ExternalID  string
	Policy      string
}

func getCredentials(c *Config) *credentials.Credentials {
	// Follow the  same priority as the AWS Terraform Provider
	// https://www.terraform.io/docs/providers/aws/#authentication

	providers := []credentials.Provider{
		&credentials.StaticProvider{Value: credentials.Value{
			AccessKeyID:     c.AccessKey,
			SecretAccessKey: c.SecretKey,
			SessionToken:    c.Token,
		}},
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{
			Filename: c.CredsFilename,
			Profile:  c.Profile,
		},
	}

	return credentials.NewChainCredentials(providers)
}

func getCredentialsFromSession(c *Config) (*credentials.Credentials, error) {
	var sess *session.Session
	var err error
	options := &session.Options{
		Config: aws.Config{
			MaxRetries: aws.Int(0),
			Region:     aws.String("us-east-1"),
		},
	}
	options.Profile = c.Profile
	options.SharedConfigState = session.SharedConfigEnable

	sess, err = session.NewSessionWithOptions(*options)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok && awsErr.Code() == "NoCredentialProviders" {
			return nil, ErrNoValidCredentialSources
		}
		return nil, fmt.Errorf("Error creating AWS session: %s", err)
	}
	creds := sess.Config.Credentials
	cp, err := sess.Config.Credentials.Get()
	if err != nil {
		return nil, ErrNoValidCredentialSources
	}

	log.Printf("[DEBUG] Got session credentials from provider: %s\n", cp.ProviderName)

	return creds, nil
}

// Client returns a properly configured ALKS client or an appropriate error if initialization fails
func (c *Config) Client() (*alks.Client, error) {
	log.Println("[DEBUG] Validting STS credentials")

	// lookup credentials
	creds := getCredentials(c)
	cp, cpErr := creds.Get()

	if cpErr == nil {
		log.Printf("[DEBUG] Got credentials from provider: %s\n", cp.ProviderName)
	}

	// validate we have credentials
	if cpErr != nil {
		if awsErr, ok := cpErr.(awserr.Error); ok && awsErr.Code() == "NoCredentialProviders" {
			var err error
			creds, err = getCredentialsFromSession(c)
			if err != nil {
				return nil, err
			}
			cp, cpErr = creds.Get()
		}
	}
	if cpErr != nil {
		return nil, ErrNoValidCredentialSources
	}

	// create a new session to test credentails
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: creds,
	})

	// validate session
	if err != nil {
		return nil, fmt.Errorf("Error creating session from STS. (%v)", err)
	}

	var stsconn *sts.STS
	// we need to assume another role before creating an ALKS client
	if c.AssumeRole.RoleARN != "" {
		arCreds := stscreds.NewCredentials(sess, c.AssumeRole.RoleARN, func(p *stscreds.AssumeRoleProvider) {
			if c.AssumeRole.SessionName != "" {
				p.RoleSessionName = c.AssumeRole.SessionName
			}

			if c.AssumeRole.ExternalID != "" {
				p.ExternalID = &c.AssumeRole.ExternalID
			}

			if c.AssumeRole.Policy != "" {
				p.Policy = &c.AssumeRole.Policy
			}
		})

		cp, cpErr = arCreds.Get()
		if cpErr != nil {
			return nil, fmt.Errorf("The role %q cannot be assumed. Please verify the role ARN, role policies and your base AWS credentials", c.AssumeRole.RoleARN)
		}

		stsconn = sts.New(sess, &aws.Config{
			Region:      aws.String("us-east-1"),
			Credentials: arCreds,
		})
	} else {
		stsconn = sts.New(sess)
	}

	// make a basic api call to test creds are valid
	cident, serr := stsconn.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	// check for valid creds
	if serr != nil {
		return nil, serr
	}

	// got good creds, create alks sts client
	client, err := alks.NewSTSClient(c.URL, cp.AccessKeyID, cp.SecretAccessKey, cp.SessionToken)

	// check if the user is using a assume-role IAM admin session or MI.
	if isValidIAM(cident.Arn, client) != true {
		return nil, errors.New("Looks like you are not using ALKS IAM credentials. This will result in errors when creating roles. \n " +
			"Note: If using ALKS CLI to get credentials, be sure to use the '-i' flag. \n Please see https://coxautoinc.sharepoint.com/sites/service-internal-tools-team/SitePages/ALKS-Terraform-Provider---Troubleshooting.aspx for more information.")
	}

	if err != nil {
		return nil, err
	}

	client.SetUserAgent(fmt.Sprintf("alks-terraform-provider-%s", getPluginVersion()))

	log.Println("[INFO] ALKS Client configured")

	return client, nil
}

func getPluginVersion() string {
	if versionNumber != "" {
		return versionNumber
	}

	return "unknown"
}

/*
	Validates ARN for assumed-role of:
		- Admin
		- IAMAdmin
		- Machine Identities.
*/
func isValidIAM(arn *string, client *alks.Client) bool {
	// Check if Admin || IAMAdmin
	if strings.Contains(*arn, "assumed-role/Admin/") || strings.Contains(*arn, "assumed-role/IAMAdmin/") || strings.Contains(*arn, "assumed-role/LabAdmin/") {
		return true
	}

	// Check if MI...
	arnParts := strings.FieldsFunc(*arn, splitBy)
	iamArn := fmt.Sprintf("arn:aws:iam::%s:role/acct-managed/%s", arnParts[3], arnParts[5])

	_, err := client.SearchRoleMachineIdentity(iamArn)
	if err != nil {
		return false
	}
	return true
}

func splitBy(r rune) bool {
	return r == ':' || r == '/'
}
