#
# PROVIDERS
#
provider "alks" {
    url      = "https://alks.coxautoinc.com/rest"
    account  = "120678615247/ALKSAdmin - awsaepnp"
    role     = "Admin"
}

# SETUP ALKS SESSION FOR AWS PROVIDER
resource "alks_session" "session" {
  name = "my_session"
  use_iam = false
}


provider "aws" {
    region     = "us-east-1"
    access_key = "${alks_session.session.access_key}"
    secret_key = "${alks_session.session.secret_key}"
    token      = "${alks_session.session.session_token}"
}

# CREATE IAM ROLE
resource "alks_iamrole" "test_role" {
    name                     = "My_Test_Role44444"
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


# CREATE SECURITY GROUP TO TEST NON-IAM
resource "aws_security_group" "btest42" {
    name   = "btest34"

    egress {
        from_port   = 0
        to_port     = 0 
        protocol    = "-1"
        cidr_blocks = [ "0.0.0.0/0" ]
    }
}
