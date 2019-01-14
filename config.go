package main

import (
	"errors"
	"fmt"
	"log"

	alks "github.com/Cox-Automotive/alks-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/ec2rolecreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

type Config struct {
	Url           string
	AccessKey     string
	SecretKey     string
	Token         string
	CredsFilename string
	Profile       string
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
		&ec2rolecreds.EC2RoleProvider{},
	}

	return credentials.NewChainCredentials(providers)
}

func (c *Config) Client() (*alks.Client, error) {
	log.Println("[DEBUG] Validting STS credentials")

	// lookup credentials
	creds := getCredentials(c)
	cp, cpErr := creds.Get()

	// validate we have credentials
	if cpErr != nil {
		return nil, errors.New(`No valid credential sources found for ALKS Provider.
Please see https://github.com/Cox-Automotive/terraform-provider-alks#authentication for more information on
providing credentials for the ALKS Provider`)
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

	// make a basic api call to test creds are valid
	stsconn := sts.New(sess)
	_, serr := stsconn.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	// check for valid creds
	if serr != nil {
		return nil, serr
	}

	// got good creds, create alks sts client
	client, err := alks.NewSTSClient(c.Url, cp.AccessKeyID, cp.SecretAccessKey, cp.SessionToken)

	if err != nil {
		return nil, err
	}

	log.Println("[INFO] ALKS Client configured")

	return client, nil
}
