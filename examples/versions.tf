terraform {
  required_providers {
    alks = {
      source  = "Cox-Automotive/alks"
      version = "1.5.8"
    }
    aws = {
      source = "hashicorp/aws"
    }
  }
  required_version = ">= 0.13"
}
