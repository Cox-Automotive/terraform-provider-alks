terraform {
  required_providers {
    alks = {
      source  = "coxautoinc.com/engineering-enablement/alks"
      version = "1.4.4"
    }
    aws = {
      source = "hashicorp/aws"
    }
  }
  required_version = ">= 0.13"
}
