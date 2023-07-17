package main

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestSuppressEquivalentTrustPolicyDiffs(t *testing.T) {
	cases := []struct {
		policy1  string
		policy2  string
		expected bool
	}{
		{
			policy1: string(`
			{
				"Version": "1234",
				"Id": "Something",
				"Statement": [
					{
						"Action": "sts:AssumeRole",
						"Effect": "Deny",
						"Sid": "2"
					},
					{
						"Action": "sts:AssumeRole",
						"Effect": "Allow",
						"Sid": "1"
					}
				]
			}
			`),
			policy2: string(`
			{
				"Id": "Something",
				"Version": "1234",
				"Statement": [
					{
						"Action": "sts:AssumeRole",
						"Effect": "Deny",
						"Sid": "2"
					},
					{
						"Action": "sts:AssumeRole",
						"Effect": "Allow",
						"Sid": "1"
					}
				]
			}
			`),
			expected: true,
		},
		{
			policy1: string(`
			{
				"Version": "2012-10-17",
				"Statement": [
				  {
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": {
					  "Service": ["lambda.amazonaws.com", "ec2.amazonaws.com"]
					}
				  }
				]
			  }
			`),
			policy2: string(`
			{
				"Version": "2012-10-17",
				"Statement": [
				  {
					"Action": "sts:AssumeRole",
					"Effect": "Allow",
					"Principal": {
					  "Service": ["ec2.amazonaws.com", "lambda.amazonaws.com"]
					}
				  }
				]
			  }
			`),
			expected: true,
		},
		{
			policy1: string(`
			{
				"Version": "1234",
				"Id": "Something",
				"Statement": [
					{
						"Action": "sts:AssumeRole",
						"Effect": "Allow",
						"Sid": "2"
					},
					{
						"Action": "sts:AssumeRole",
						"Effect": "Allow",
						"Sid": "1"
					}
				]
			}
			`),
			policy2: string(`
			{
				"Id": "Something",
				"Version": "1234",
				"Statement": [
					{
						"Action": "sts:AssumeRole",
						"Effect": "Deny",
						"Sid": "2"
					},
					{
						"Action": "sts:AssumeRole",
						"Effect": "Allow",
						"Sid": "1"
					}
				]
			}
			`),
			expected: false,
		},
	}
	for _, c := range cases {
		isEqual := SuppressEquivalentTrustPolicyDiffs("key", c.policy1, c.policy2, &schema.ResourceData{})
		if isEqual != c.expected {
			t.Fatalf("Error matching output and expected: %#v vs %#v", isEqual, c.expected)
		}
	}
}
