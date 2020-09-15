#
# PROVIDERS
#
provider "alks" {
  url = "https://alks.foo.com/rest"
}

provider "alks" {
  url     = "https://alks.foo.com/rest"
  account = "<account No>"
  role    = "<role>"
  alias   = "second"
}


provider "aws" {
  region = "us-east-1"
}

# CREATE IAM ROLE -- Initial Provider
resource "alks_iamrole" "test_role" {
  name                     = "TEST-DELETE"
  type                     = "AWS CodeBuild"
  include_default_policies = false
  enable_alks_access       = true
}

# CREATE IAM ROLE -- Secondary Provider
resource "alks_iamrole" "test_role_nonprod" {
  provider                 = alks.second
  name                     = "TEST-DELETE"
  type                     = "AWS CodeBuild"
  include_default_policies = false
  enable_alks_access       = true
}

# ATTACH POLICY
resource "aws_iam_role_policy" "test_policy" {
    name     = "test_policy"
    role     = "${alks_iamrole.test_role.name}"
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

# ATTACH MANAGED POLICY
resource "aws_iam_role_policy_attachment" "sr-attach" {
    role       = "${alks_iamrole.test_role.name}"
    policy_arn = "arn:aws:iam::aws:policy/service-role/AWSElasticBeanstalkService"
}

# CREATE LTK USER
resource "alks_ltk" "ltk" {
    iam_username = "TEST_LTK_USER"
}