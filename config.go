package main

import (
	"fmt"
	"github.com/Cox-Automotive/alks-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"log"
	"strings"
)

type Config struct {
	Url           string
	AccessKey     string
	SecretKey     string
	Token         string
	CredsFilename string
	Profile       string
	Account       string
}

func parseAccountInfoFromArn(arn string) (string, string, error) {
	parts := strings.Split(arn, ":")
	if len(parts) < 5 {
		return "", "", fmt.Errorf("Unable to parse ID from invalid ARN: %v", arn)
	}

	return parts[1], parts[4], nil
}

func getCredentials(c *Config) (*credentials.Credentials, error) {
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

	return credentials.NewChainCredentials(providers), nil
}

func (c *Config) Client() (*alks.Client, error) {
	log.Println("[DEBUG] Validting STS credentials")

	creds, cErr := getCredentials(c)

	if cErr != nil {
		return nil, cErr
	}

	cp, cpErr := creds.Get()

	if cpErr != nil {
		return nil, cpErr
	}

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: creds,
	})

	if err != nil {
		return nil, fmt.Errorf("Error creating session from STS. (%v)", err)
	}

	stsconn := sts.New(sess)
	outCallerIdentity, err := stsconn.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		return nil, err
	}

	_, accountId, err := parseAccountInfoFromArn(*outCallerIdentity.Arn)

	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(c.Account, accountId) {
		return nil, fmt.Errorf("The provided STS token is not valid for the provided account (%v). It's for %v.", c.Account, accountId)
	}

	// now we know we have a valid STS
	client, err := alks.NewSTSClient(c.Url, cp.AccessKeyID, cp.SecretAccessKey, cp.SessionToken, c.Account)

	if err != nil {
		return nil, err
	}

	log.Println("[INFO] ALKS Client configured")

	return client, nil
}
