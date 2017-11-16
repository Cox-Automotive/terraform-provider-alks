#
# PROVIDERS
#
provider "alks" {
    url      = "https://alks.foo.com/rest"
    account  = "124124124124124/ALKSAdmin - awsfoo420labs"
}

provider "aws" {
    region     = "us-east-1"
}

# CREATE IAM ROLE
resource "alks_iamrole" "test_role" {
    name                     = "aba-test-123456"
    type                     = "Amazon EC2"
    include_default_policies = false
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