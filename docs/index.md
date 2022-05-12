# ALKS Terraform Provider

The ALKS Terraform Provider is used to interact with several resources supported by ALKS. The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the ALKS Terraform Provider
provider "alks" {
    url     = "https://alks.foo.com/rest"
    version = ">= 2.2.0"
}

# Create an ALKS IAM role
resource "alks_iamrole" "test_role" {
    name                     = "My_Test_Role"
    type                     = "Amazon EC2"
    include_default_policies = false
    enable_alks_access       = false
}
```

## Authentication

The ALKS Terraform Provider offers a flexible means of providing credentials for authentication. The following methods are supported, in this order, and explained below:

### Static credentials

!> **Warning:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system.

Static credentials can be provided via an `access_key`, `secret_key` and `token` in-line in the ALKS provider block. This method is generally not recommended, since the credentials could accidentally be committed or shared.

```hcl
provider "alks" {
    url        = "https://alks.foo.com/rest"
    version    = ">= 2.2.0"
    access_key = "accesskey"
    secret_key = "secretkey"
    token      = "sessiontoken"
}
```

### Environment variables
You can provide your credentials via the `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_SESSION_TOKEN` environment variables. If you need to pass multiple AWS credentials (when using a combination of Roles, like PowerUser and IAMAdmin) you can use the ALKS_ prefix in place of AWS_ (ex: ALKS_ACCESS_KEY_ID) as these are prioritized over the AWS_ prefixed environment varaibles for the provider.

Terraform file:
```hcl
provider "alks" {
    url     = "https://alks.foo.com/rest"
    version = ">= 2.2.0"
}
```

Terminal:
```hcl
$ alks sessions open -i
$ export AWS_ACCESS_KEY_ID="accesskey"
$ export AWS_SECRET_ACCESS_KEY="secretkey"
$ export AWS_SESSION_TOKEN="sessiontoken"
$ terraform plan
```

### Shared credentials/configuration file
You can use an AWS credentials file to specify your credentials. The default location is `$HOME/.aws/credentials` on Linux and OSX, or `"%USERPROFILE%\.aws\credentials"` for Windows users. If we fail to detect credentials inline, or in the environment, Terraform will check this location last. You can optionally specify a different location in the configuration via the `shared_credentials_file` attribute, or via the environment with the `AWS_SHARED_CREDENTIALS_FILE` variable. This method also supports a profile configuration and matching `AWS_PROFILE` environment variable.

```hcl
provider "alks" {
    url                     = "https://alks.foo.com/rest"
    version                 = ">= 2.2.0"
    shared_credentials_file = "/Users/my_user/.aws/credentials"
    profile                 = "foo"
}
```

### Machine Identities
You can use a role created with ALKS with the `enable_alks_access` flag set to `true` to authenticate requests against ALKS.

In order to do this, ALKS must be called from within AWS using STS credentials from an instance profile associated with the role with `enable_alks_access` set. This also works from Lambda functions in the same way.

The STS credentials are used and provided in the same way that the AWS CLI uses the credentials, so there is nothing special you have to do to use Machine Identities.

Your ALKS provider block can look just like this:

```hcl
provider "alks" {
    url     = "https://alks.foo.com/rest"
    version = ">= 2.2.0"
}
```

Since Machine Identities work with Instance Profile Metadata directly, it can be helpful to assume another role or cross account trust. For example:

```hcl
provider "alks" {
   url     = "https://alks.foo.com/rest"
   version = ">= 2.2.0"
   assume_role {
      role_arn = "arn:aws:iam::112233445566:role/acct-managed/JenkinsPRODAccountTrust"
   }
}
```


### Multiple Provider Configuration

You can configure multiple ALKS providers to each have their own account context. 

The initial provider must have credentials set in a default way (static, shared credentials file, environment variables, etc) before the second provider can determine whether your account/role combination are allowed. 

The second (or so) provider can then be used to generate resources for multiple accounts in one plan / apply.

Note: This only works for accounts you have access to!

```tf
# PROVIDER 1
provider "alks" {
  url = "https://alks.coxautoinc.com/rest"
}

# PROVIDER 2
provider "alks" {
  url     = "https://alks.coxautoinc.com/rest"
  account = "<account No>"
  role    = "<role>"
  alias   = "second"
}

# CREATE IAM ROLE -- PROVIDER 1
resource "alks_iamrole" "test_role" {
  name                     = "TEST-DELETE"
  type                     = "AWS CodeBuild"
  include_default_policies = false
  enable_alks_access       = true
}

# CREATE IAM ROLE -- PROVIDER 2
resource "alks_iamrole" "test_role_nonprod" {
  provider                 = alks.second
  name                     = "TEST-DELETE"
  type                     = "AWS CodeBuild"
  include_default_policies = false
  enable_alks_access       = true
}
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html?_ga=2.182283811.562816692.1597670778-20010454.1565803281) (e.g. `alias` and `version`), the following arguments are supported in the AWS provider block:

* `url` - (Required) The URL to your ALKS server. Also read from ENV.ALKS_URL
* `access_key` - (Optional) The access key from a valid STS session. Also read from ENV.ALKS_ACCESS_KEY_ID and ENV.AWS_ACCESS_KEY_ID.
* `secret_key` - (Optional) The secret key from a valid STS session. Also read from ENV.ALKS_SECRET_ACCESS_KEY and ENV.AWS_SECRET_ACCESS_KEY.
* `token` - (Optional) The session token from a valid STS session. Also read from ENV.ALKS_SESSION_TOKEN and ENV.AWS_SESSION_TOKEN.
* `shared_credentials_file` - (Optional) The the path to the shared credentials file. Also read from ENV.AWS_SHARED_CREDENTIALS_FILE.
* `profile` - (Optional) This is the AWS profile name as set in the shared credentials file. Also read from ENV.AWS_PROFILE.
* `assume_role` - (Optional) This is the role information to assume before making calling ALKS. This feature works the same as the assume_role feature of the AWS Terraform Provider.
    * `role_arn` - (Required) The Role ARN to assume for calling the ALKS API.
    * `session_name` - (Optional) The session name to provide to AWS when creating STS credentials. Please see the AWS SDK documentation for more information.
    * `external_id` - (Optional) The external identifier to provide to AWS when creating STS credentials. Please see the AWS SDK documentation for more information.
    * `policy` - (Optional) This specifies additional policy restrictions to apply to the resulting STS credentials beyond any existing inline or managed policies. Please see the AWS SDK documentation for more information.

---
### Supported Versions

| Terraform 0.10.x       |  Terraform 0.11.x        | Terraform 0.12.x - Terraform 0.15.x  | Terraform 1.0.x   |
| ---------------------- | ------------------------ | ------------------------------------ | ----------------- |
| ALKS TFP 0.9.0 < 1.3.0 | ALKS TFP 1.3.0 < 1.5.15  | ALKS TFP 1.3.0+                      | ALKS TFP 2.0.0+   |

For questions, please reach out to the [ALKS team](https://github.com/orgs/Cox-Automotive/teams/cai-internal-tools).