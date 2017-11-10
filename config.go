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
	Url       string
	AccessKey string
	SecretKey string
	Token     string
	Account   string
}

func parseAccountInfoFromArn(arn string) (string, string, error) {
	parts := strings.Split(arn, ":")
	if len(parts) < 5 {
		return "", "", fmt.Errorf("Unable to parse ID from invalid ARN: %v", arn)
	}

	return parts[1], parts[4], nil
}

func (c *Config) Client() (*alks.Client, error) {
	log.Println("[DEBUG] Validting STS credentials")

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String("us-east-1"),
		Credentials: credentials.NewStaticCredentials(c.AccessKey, c.SecretKey, c.Token),
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
	client, err := alks.NewSTSClient(c.Url, c.AccessKey, c.SecretKey, c.Token, c.Account)

	if err != nil {
		return nil, err
	}

	log.Println("[INFO] ALKS Client configured")

	return client, nil
}
