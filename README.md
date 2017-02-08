# ALKS Provider for Terraform

This module is used for creating IAM Roles via the ALKS API. This is intended to be used by user's with the `Admin` role, if you have an `IAMAdmin` role you'll need to have it upgraded to `Admin`.

## Building & Testing

- Install Dependencies: `make deps`
- Build: `make build`
- TF Plan: `make plan`
- Test: `make test`

## Example

First you'll need to export a valid ALKS Admin session to your environment variables. Be sure to insert your own account, this will be used for creating all AWS resources with the exception of roles since those will use the alks role.

`eval $(alks sessions open -i -a "######/ALKSAdmin - awsaepnp" -r "Admin")`

When prompted enter your network username and password, or export them as environment variables. (`ALKS_USERNAME` and `ALKS_PASSWORD`). You'll also need to provide the same account and role that you provided to `alks sessions` above in your alks provider definition.

Then simply run: `terraform plan` and then `terraform apply`

## Common Build Issues

The following build issue can occur if you have an older version of the terraform executable installed.

```
Error configuring: 1 error(s) occurred:
* Incompatible API version with plugin. Plugin version: 3, Ours: 2
```

To build the provider plugin against a specific terraform version, checkout the version in the $GOPATH/src/github.com/hashicorp/terraform source directory and rebuild. e.g.

```
$ terraform --version
Terraform v0.8.6

$ ( cd $GOPATH/src/github.com/hashicorp/terraform && git checkout v0.8.6 )
```

To revert back to building againt the latest source change the git checkout statement to git checkout master