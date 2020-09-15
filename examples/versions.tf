terraform {
  required_providers {
    alks = {
      source = "coxautoinc.com/engineering-enablement/alks"
    }
    aws = {
      source = "hashicorp/aws"
    }
  }
  required_version = ">= 0.13"
}
