#
# PROVIDERS
#
provider "aws" {
    alias                    = "main"
    region                   = "us-east-1"
    shared_credentials_file  = "/Users/brianantonelli/.aws/credentials"
    profile                  = "tf-main"
}

provider "aws" {
    alias                    = "iam"
    region                   = "us-east-1"
    shared_credentials_file  = "/Users/brianantonelli/.aws/credentials"
    profile                  = "tf-iam"
}

provider "alks" {
    url      = ""
    account  = ""
    role     = ""
    username = ""
    password = ""
}

# CREATE IAM ROLE
resource "alks_iamrole" "mah_role" {
    name                     = "My_Test_Role5"
    type                     = "Amazon EC2"
    include_default_policies = false
}

# ATTACH POLICY
resource "aws_iam_role_policy" "test_policy" {
    provider = "aws.iam"
    name     = "test_policy"
    role     = "${alks_iamrole.mah_role.name}"
    policy   = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}