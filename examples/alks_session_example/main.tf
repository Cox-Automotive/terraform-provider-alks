terraform {
  required_version = "~> 1.5"

  required_providers {
    alks = {
      source  = "cox-automotive/alks"
      #source  = "terraformLocal/alks"
      #source = "/Users/james.barcelo/code/github.com/jcarlson/terraform-provider-alks"
      #version = "~> 2.8"
    }

    shell = {
      source  = "scottwinkler/shell"
      version = "~> 1.7"
    }
  }
}

provider "alks" {
  url = "https://alks.coxautoinc.com/rest"

  ignore_tags {
    key_prefixes = ["cai:catalog"]
  }
}

data "alks_keys" "session" {
}

data "alks_session" "current" {
}

resource "shell_script" "sts_identity" {
  environment = {
    AWS_REGION = "us-east-1"
  }

  sensitive_environment = {
    AWS_ACCESS_KEY_ID     = data.alks_session.current.access_key
    AWS_SECRET_ACCESS_KEY = data.alks_session.current.secret_key
    AWS_SESSION_TOKEN     = data.alks_session.current.session_token
  }

  lifecycle_commands {
    create = "aws sts get-caller-identity"
    delete = ""
  }

  lifecycle {
    ignore_changes = [
      sensitive_environment["AWS_ACCESS_KEY_ID"],
      sensitive_environment["AWS_SECRET_ACCESS_KEY"],
      sensitive_environment["AWS_SESSION_TOKEN"]
    ]
  }
}

output "user_id" {
  value = shell_script.sts_identity.output["UserId"]
}

output "arn" {
  value = shell_script.sts_identity.output["Arn"]
}

output "alks_session_access_key" {
  value = data.alks_session.current.access_key
}

output "alks_session_secret_key" {
  value = data.alks_session.current.secret_key
  sensitive = true

}

output "alks_session_session_token" {
  value = data.alks_session.current.session_token
  sensitive = true
}

output "alks_keys_access_key" {
  value = data.alks_keys.session.access_key
}

output "alks_keys_secret_key" {
  value = data.alks_keys.session.secret_key
}
  
