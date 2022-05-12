ALKS Provider for Terraform
=========

[![Build Status](https://travis-ci.org/Cox-Automotive/terraform-provider-alks.svg?branch=master)](https://travis-ci.org/Cox-Automotive/terraform-provider-alks)

This provider is for creating and modifying IAM roles via the ALKS API.

## Pre-Requisites

* An ALKS Admin or IAMAdmin STS [assume-role](http://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html) session needed. PowerUser access is not sufficient to create IAM roles.
    * This tool is best suited for users with an `Admin` role
    * With an `IAMAdmin|LabAdmin` role, you can create roles and attach policies, but you can't create other infrastructure.
* Works with [Terraform](https://www.terraform.io/) version `0.10.0` or newer.
    * [ALKS TFP - Supported Versions Matrix](https://registry.terraform.io/providers/Cox-Automotive/alks/latest/docs#supported-versions)

## Building from Source

To build the ALKS provider, install [Go](http://www.golang.org/) (preferably version 1.14.4 or greater).

Clone this repository and `cd` into the cloned directory. All the necessary dependencies are vendored, so type `make build test` to build and test the project. If this exits with exit status `0`, then everything is working! Check your `examples` directory for an example Terraform script and the generated binary.

```bash
git clone https://github.com/Cox-Automotive/terraform-provider-alks.git
cd terraform-provider-alks
make build test
```

As stated above, if the tests run, everything is working. What you won't see are any passing tests. To get to an operational testing state:
1. set the ALKS_URL and TF_ACC environment variables
```bash
export TF_ACC=true
export ALKS_URL=https://dev.alks.coxautoinc.com/rest
```
2. Copy the environment variables from **CoxAT Labs 95 (ALKS Dev)** into your terminal
```bash
export AWS_ACCESS_KEY_ID=<key_from_alks_web>
export AWS_SECRET_ACCESS_KEY=<key_from_alks_web>
export AWS_SESSION_TOKEN=<token_from_alks_web>
export AWS_DEFAULT_REGION=us-east-1
```
If an error stating `Role already exists with the same name: <role-name>` is encountered during testing (errored out tests do not initiate resource tear down), navigate to the AWS console of the Labs account from ALKSWeb and manually delete the IAM role listed in the error.

If you need any additional dependencies while developing, add the dependency by running `go get <dependency>` and then add it to the vendor folder by running `go mod vendor`.

## Updating The Version of alks-go
If using VSCode, hover over the versioned alks-go import in go.mod and click on the link to go.dev (Go package index). The latest version should have the commit hash at the HEAD of master (you may have to wait for the site to update ~ 20 min). Copy the version number and paste it over the previous version in go.mod. On the command line:
```bash
go mod download github.com/Cox-Automotive/alks-go
go mod vendor
go mod tidy
```
At this point, the dependency should reflect the state of alks-go's master branch
## Documentation

Documentation and examples can be found on the [Terraform website](https://registry.terraform.io/providers/Cox-Automotive/alks/latest/docs).
