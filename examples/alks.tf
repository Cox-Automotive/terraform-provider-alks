#
# PROVIDERS
#
provider "aws" {
    region = "us-east-1"
}

provider "alks" {
    url      = "https://alks.foobar.baz/rest"
    account  = "2352352352/ALKSAdmin - awslabs"
    role     = "Admin"
}

# CREATE IAM ROLE
resource "alks_iamrole" "test_role" {
    name                     = "My_Test_Role4444"
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
    name   = "btest3"

    egress {
        from_port   = 0
        to_port     = 0 
        protocol    = "-1"
        cidr_blocks = [ "0.0.0.0/0" ]
    }
}
