ALKS Provider for Terraform
=========

[![Build Status](https://travis-ci.org/Cox-Automotive/terraform-provider-alks.svg?branch=master)](https://travis-ci.org/Cox-Automotive/terraform-provider-alks)

This provider is for creating and modifying IAM roles via the ALKS API.

## Pre-Requisites

* An ALKS Admin or IAMAdmin STS [assume-role](http://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html) session is needed. PowerUser access is not sufficient to create IAM roles.
    * This tool is best suited for users with an `Admin` role
    * With an `IAMAdmin|LabAdmin` role, you can create roles and attach policies, but you can't create other infrastructure.
* Works with [Terraform](https://www.terraform.io/) version `0.10.0` or newer.

## Building from Source

To build the ALKS provider, install [Go](http://www.golang.org/) (preferably version 1.14.4 or greater).

Clone this repository and `cd` into the cloned directory. All the necessary dependencies are vendored, so type `make build test` to build and test the project. If this exits with exit status `0`, then everything is working! Check your `examples` directory for an example Terraform script and the generated binary.

```bash
git clone https://github.com/Cox-Automotive/terraform-provider-alks.git
cd terraform-provider-alks
make build test
```

If you need any additional dependencies while developing, add the dependency by running `go get <dependency>` and then add it to the vendor folder by running `go mod vendor`.

## Documentation

Documentation and examples can be found on the [Terraform website](https://registry.terraform.io/providers/Cox-Automotive/alks/latest/docs).