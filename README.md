ALKS Provider for Terraform
=========

[![Build Status](https://travis-ci.org/Cox-Automotive/terraform-provider-alks.svg?branch=master)](https://travis-ci.org/Cox-Automotive/terraform-provider-alks)

This module is used for creating IAM Roles via the ALKS API.

## Pre-Requisites

* This is intended to be used by user's with the `Admin` role, if you have an `IAMAdmin` role you'll be able to create roles and attach policies, but you won't be able to create other infrastructure. Upgrading to `Admin` will remedy this.
* You'll need to make sure you've [installed Terraform](https://www.terraform.io/) version 0.7 or newer.

## Installation

1. Download the correct binary for your platform from [Releases](https://github.com/Cox-Automotive/terraform-provider-alks/releases) and place it somewhere on your filesystem.
2. Configure Terraform to find the plugin by creating `~/.terraformrc` on *nix and `%APPDATA%/terraform.rc` for Windows.

### Example Provider Config

If you didn't add terraform-provider-alks to your path, you'll need to put the full path to the location of the plugin.

```
providers {
    alks = "/path/to/terraform-provider-alks-<platform>-<arch>"
}
```

## Usage

First you'll need to export a valid ALKS session to your environment variables. Be sure to insert your own account, this will be used for attaching policies as well as creating all AWS resources. The alks provider is only responsible for creating the initial role.

`eval $(alks sessions open -i -a "######/ALKSAdmin - sdgsgasf" -r "Admin")`

### Provider Configuration

#### `alks`

```
provider "alks" {
    url      = "<ALKS_URL>"
    account  = "<ALKS_ACCOUNT>"
    role     = "<ALKS_ROLE>"
    username = "<ALKS_USERNAME>"
    password = "<ALKS_PASSWORD>"
}
```

Provider Options:
* `url` - (Required) The URL to your ALKS server. Also read from `ENV.ALKS_URL`
* `account` - (Required) The ALKS account to use. Also read from `ENV.ALKS_ACCOUNT`
* `role` - (Required) The ALKS role to use. This should be `Admin` or `IAMAdmin` Also read from `ENV.ALKS_ROLE`
* `username` - (Required) The username you use to login to ALKS. Read from `ENV.ALKS_USERNAME` - **Should be provided via env vars and not stored in your TF files.**
* `password` - (Required) The password you use to login to ALKS. Also read from `ENV.ALKS_PASSWORD` - **Should be provided via env vars and not stored in your TF files.**

All of these options should match what you configured with the ALKS CLI.

### Resource Configuration

#### `alks_iamrole`

```
resource "alks_iamrole" "test_role" {
    name                     = "My_Test_Role"
    type                     = "Amazon EC2"
    include_default_policies = false
}
```

Value                             | Type     | Forces New | Value Type | Description
--------------------------------- | -------- | ---------- | ---------- | -----------
`name`                           | Required | yes        | string     | The name of the IAM role to create. This parameter allows a string of characters consisting of upper and lowercase alphanumeric characters with no spaces. You can also include any of the following characters: =,.@-. Role names are not distinguished by case.
`type`                           | Required | yes        | string     | The role type to use. [Available Roles](https://gist.github.com/brianantonelli/5769deff6fd8f3ff30e40b844f0b1fb4)
`include_default_policies`                           | Required | yes        | bool     | Whether or not the default managed policies should be attached to the role.
`role_added_to_ip`                           | Computed | n/a        | bool     | Indicates whether or not an instance profile role was created.
`arn`                           | Computed | n/a        | string     | Provides the ARN of the role that was created.
`ip_arn`                           | Computed | n/a        | string     | If `role_added_to_ip` was `true` this will provide the ARN of the instance profile role.

## Example

Check out `test.tf` for an very basic Terraform script which:

1. Creates an AWS provider and ALKS provider
2. Creates an IAM Role via the ALKS provider
3. Attaches a policy to the created role using the AWS provider
4. Creates a security group using the AWS provider

This example is meant to show how you would combine a typical AWS Terraform script with our custom provider in order to automate the creation of IAM roles.

## Building from Source

- Install Dependencies: `make get-deps`
- Build: `make build`
- TF Plan: `make plan`
- Test: `make test`

### Common Build Issues

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
