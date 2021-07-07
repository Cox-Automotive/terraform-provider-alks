#
# PROVIDERS
#
provider "alks" {
  url = "https://alks.foo.com/rest"
}

# Second ALKS provider, for an account I have access to.
provider "alks" {
  url     = "https://alks.foo.com/rest"
  account = "<account No>"
  role    = "<role>"
  alias   = "nonprod"
}

data "alks_keys" "non_prod_keys" {
  provider = alks.nonprod
}

provider "aws" {
  region = "us-east-1"
}

# Second AWS provider, using credentials retreived from data source.
provider "aws" {
  region = "us-east-1"
  alias  = "nonprod"

  # data source alks keys
  access_key = data.alks_keys.non_prod_keys.access_key
  secret_key = data.alks_keys.non_prod_keys.secret_key
  token      = data.alks_keys.non_prod_keys.session_token
}

# CREATE IAM ROLE -- Initial Provider
resource "alks_iamrole" "test_role" {
  name                     = "TEST-DELETE-PRIMARY-PROVIDER"
  type                     = "AWS CodeBuild"
  include_default_policies = false
  enable_alks_access       = true
}

resource "alks_iamrole" "test_dynamic_role" {
  name                     = "TEST-DELETE-DYNAMIC"
  type                     = "Amazon EKS IRSA"
  include_default_policies = false
  enable_alks_access       = true
  template_fields = {
    "OIDC_PROVIDER"       = "abc1234"
    "K8S_NAMESPACE"       = "default"
    "K8S_SERVICE_ACCOUNT" = "definitely-not-real"
  }
}

# CREATE IAM ROLE -- Secondary Provider
resource "alks_iamrole" "test_role_nonprod" {
  provider                 = alks.nonprod
  name                     = "TEST-DELETE-SECONDARY-PROVIDER"
  type                     = "AWS CodeBuild"
  include_default_policies = false
  enable_alks_access       = true
}

# ATTACH POLICY
resource "aws_iam_role_policy" "test_policy" {
  name   = "test_policy"
  role   = alks_iamrole.test_role.name
  policy = <<EOF
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
  role       = alks_iamrole.test_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AWSElasticBeanstalkService"
}

# CREATE LTK USER
resource "alks_ltk" "ltk" {
  iam_username = "TEST-LTK-USER"
}
