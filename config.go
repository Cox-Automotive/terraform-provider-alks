package main

import (
	"fmt"
	"github.com/Cox-Automotive/alks-go"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"log"
)

type Config struct {
	Url       string
	AccessKey string
	SecretKey string
	Token     string
	Account   string
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
	_, err = stsconn.GetCallerIdentity(&sts.GetCallerIdentityInput{})

	if err != nil {
		return nil, err
	}

	// TODO: validate account matches?

	// now we know we have a valid STS
	client, err := alks.NewSTSClient(c.Url, c.AccessKey, c.SecretKey, c.Token, c.Account)

	if err != nil {
		return nil, err
	}

	log.Println("[INFO] ALKS Client configured")

	return client, nil
}
