terraform {
  required_providers {
    alks = {
      source  = "Cox-Automotive/alks"
      version = "2.0.4"
    }
    aws = {
      source = "hashicorp/aws"
    }
  }
  required_version = ">= 0.13"
}
