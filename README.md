ALKS Provider for Terraform
=========

This module is used for creating IAM Roles via the ALKS API.

[![Build Status](https://travis-ci.org/Cox-Automotive/terraform-provider-alks.svg?branch=master)](https://travis-ci.org/Cox-Automotive/terraform-provider-alks)

## Pre-Requisites

This is intended to be used by user's with the `Admin` role, if you have an `IAMAdmin` role you'll be able to create roles and attach policies, but you wont't be able to create other infrastructure. Upgrading to `Admin` will remedy this.

## Example

Check out `test.tf` for an very basic Terraform script which:

1. Creates an AWS provider and ALKS provider
2. Creates an IAM Role via the ALKS provider
3. Attaches a policy to the created role using the AWS provider
4. Creates a security group using the AWS provider

This example is meant to show how you would combine a typical AWS Terraform script with our custom provider in order to automate the creation of IAM roles.

First you'll need to export a valid ALKS session to your environment variables. Be sure to insert your own account, this will be used for attaching policies as well as creating all AWS resources. The alks provider is only responsible for creating the role.

`eval $(alks sessions open -i -a "######/ALKSAdmin - sdgsgasf" -r "Admin")`

When prompted enter your network username and password, or export them as environment variables. (`ALKS_USERNAME` and `ALKS_PASSWORD`). You'll also need to provide the same account and role that you provided to `alks sessions` above in your alks provider definition.

Then simply run: `terraform plan` and then `terraform apply`

## Building from Source

- Install Dependencies: `make get-deps`
- Build: `make build`
- TF Plan: `make plan`
- Test: `make test`

## Common Build Issues

The following build issue can occur if you have an older version of the terraform executable installed.

```
Error configuring: 1 error(s) occurred:
* Incompatible API version with plugin. Plugin version: 3, Ours: 2
```

To build the provider plugin against a specific terraform version, checkout the version in the `$GOPATH/src/github.com/hashicorp/terraform` source directory and rebuild.

```
$ terraform --version
Terraform v0.8.6

$ ( cd $GOPATH/src/github.com/hashicorp/terraform && git checkout v0.8.6 )
```

To revert back to building againt the latest source change the git checkout statement to `git checkout master`.