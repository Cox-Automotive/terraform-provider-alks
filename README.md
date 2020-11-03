ALKS Provider for Terraform
=========

[![Build Status](https://travis-ci.org/Cox-Automotive/terraform-provider-alks.svg?branch=master)](https://travis-ci.org/Cox-Automotive/terraform-provider-alks)

This provider is for creating and modifying IAM roles via the ALKS API.

## Pre-Requisites

* An ALKS Admin or IAMAdmin STS [assume-role](http://docs.aws.amazon.com/STS/latest/APIReference/API_AssumeRole.html) session is needed. PowerUser access is not sufficient to create IAM roles.
    * This tool is best suited for users with an `Admin` role
    * With an `IAMAdmin|LabAdmin` role, you can create roles and attach policies, but you can't create other infrastructure.
* Works with [Terraform](https://www.terraform.io/) version `0.10.0` or newer.


## Local Installation

### Terraform version < 0.13 Installation

* Download and install [Terraform](https://www.terraform.io/intro/getting-started/install.html)

* Download ALKS Provider binary for your platform from [Releases](https://github.com/Cox-Automotive/terraform-provider-alks/releases)

For example on macOS:

```
curl https://github.com/Cox-Automotive/terraform-provider-alks/releases/download/1.5.0/terraform-provider-alks_1.5.0_darwin_amd64.zip -O -J -L | unzip
```

* Configure Terraform to use this plugin by placing the binary in `.terraform.d/plugins/` on MacOS/Linux or `terraform.d\plugins\` in your user's "Application Data" directory on Windows.

* Note: If you've used a previous version of the ALKS provider and created a `.terraformrc` file in your home directory you'll want to remove it prior to updating.

### Terraform version >= 0.13 Terraform Installation

* Download and install [Terraform](https://www.terraform.io/intro/getting-started/install.html)

* Download ALKS Provider binary for your platform from [Releases](https://github.com/Cox-Automotive/terraform-provider-alks/releases)
  
For example on macOS:

```
curl https://github.com/Cox-Automotive/terraform-provider-alks/releases/download/1.5.0/terraform-provider-alks_1.5.0_darwin_amd64.zip -O -J -L | unzip
```

* Go into the Terraform plugins path; `.terraform.d/plugins/` on MacOS/Linux or `terraform.d\plugins\` in your user's "Application Data" directory on Windows.

* Create the following directories: `coxautoinc.com/engineering-enablement/alks/1.5.0/<OS>_<ARCH>` and put the binary into the `<OS>_<ARCH>/` directory.
  * Note: This `<OS>_<ARCH>` will vary depending on your system. For example, 64-bit MacOS would be: `darwin_amd64` while 64-bit Windows 10 would be: `windows_amd64` 

* Finally, configure Terraform.
    * In your `versions.tf` or `main.tf` file you'll want to add the new ALKS provider as such:
    ```
    terraform {
        required_version = ">= 0.13"
        required_providers {
        alks = {
            source = "coxautoinc.com/engineering-enablement/alks"
            }
        }
    }
    ```

* Note: If you've previously installed our provider and it is stored in your remote state, you may need to run the [`replace-provider` command](https://www.terraform.io/docs/commands/state/replace-provider.html).

## Example

See [this example](examples/alks.tf) for a basic Terraform script which:

1. Creates an AWS provider and ALKS provider
2. Creates an IAM Role via the ALKS provider
3. Attaches a policy to the created role using the AWS provider

This example is intended to show how to combine a typical AWS Terraform script with the ALKS provider to automate the creation of IAM roles and other infrastructure.

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

Documentation can be found on the [Terraform website](https://registry.terraform.io/providers/Cox-Automotive/alks/latest/docs).