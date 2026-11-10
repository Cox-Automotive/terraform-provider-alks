# ALKS Terraform Provider

The ALKS Terraform Provider is used to interact with several resources supported by ALKS. The provider needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the ALKS Terraform Provider
provider "alks" {
    url     = "https://alks.foo.com/rest"
}

# Create an ALKS IAM role
resource "alks_iamrole" "test_role" {
    name                     = "My_Test_Role"
    assume_role_policy       = jsonencode({
        Version = "2012-10-17",
        Statement = [
            {
                Action = "sts:AssumeRole",
                Effect = "Allow",
                Principal = {
                    Service = "ec2.amazonaws.com"
                },
                Sid = ""
            }
        ]
    })
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
    version    = ">= 2.6.0"
    access_key = "accesskey"
    secret_key = "secretkey"
    token      = "sessiontoken"
}
```

### Environment variables
You can provide your credentials via the `AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY` and `AWS_SESSION_TOKEN` environment variables. The `ALKS_` names (for example `ALKS_ACCESS_KEY_ID`) are also read, but only when the matching `AWS_` name is not set: `AWS_*` takes precedence. If you need the ALKS provider to use different credentials from other AWS tooling in the same shell (for example PowerUser for the AWS provider and IAMAdmin for ALKS), use the `ALKSTFP_` names (for example `ALKSTFP_ACCESS_KEY_ID`), which take precedence over `AWS_*` for this provider only. See [Environment variables](#environment-variables-reference) for every variable.

~> **Changed in 2.9.0:** earlier versions of this page said `ALKS_` credentials take priority over `AWS_`. The provider has always used `AWS_` first; the page was wrong, and the behaviour is unchanged.

Terraform file:
```hcl
provider "alks" {
    url     = "https://alks.foo.com/rest"
    version = ">= 2.6.0"
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
    version                 = ">= 2.6.0"
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
    version = ">= 2.6.0"
}
```

Since Machine Identities work with Instance Profile Metadata directly, it can be helpful to assume another role or cross account trust. For example:

```hcl
provider "alks" {
   url     = "https://alks.foo.com/rest"
   version = ">= 2.6.0"
   assume_role {
      role_arn = "arn:aws:iam::112233445566:role/acct-managed/JenkinsPRODAccountTrust"
   }
}
```

### Tags
You can specify tags to add to all of your roles created with ALKS by using the `default_tags` block in the provider configuration.  You can also choose to ignore existing tags on a resource by including tag keys or key prefixes in the `ignore_tags` block.  These ignored tags will not show up on Terraform Plans or Applys, and will not be removed from the resource by Terraform. 

You may also specify tags on individual roles using the `tags` block.

Your ALKS configuration could look like this:

```hcl
provider "alks" {
  url     = "https://alks.foo.com/rest"
  version = ">= 2.3.0"
  default_tags {
    tags = {
      "defaultTagKey" = "defaultTagValue"
    }
  }
  ignore_tags {
    keys         = ["ignoreThisKey"]
    key_prefixes = ["cai:", "coxauto:"]
  }
}
resource "alks_iamrole" "test_role" {
  name                     = "My_Test_Role"
  type                     = "Amazon EC2"
  include_default_policies = false
  enable_alks_access       = false
  tags = {
    "roleSpecificTagKey" = "value"
  }
}
```

Note: Role specific tag values will overwrite default values if the key appears in both places.


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

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html?_ga=2.182283811.562816692.1597670778-20010454.1565803281) (e.g. `alias` and `version`), the following arguments are supported in the ALKS provider block:

* `url` - (Required) The URL to your ALKS server. Also read from ENV.ALKSTFP_URL or ENV.ALKS_URL
* `access_key` - (Optional) The access key from a valid STS session. Also read from ENV.ALKSTFP_ACCESS_KEY_ID, ENV.AWS_ACCESS_KEY_ID and ENV.ALKS_ACCESS_KEY_ID, in that order.
* `secret_key` - (Optional) The secret key from a valid STS session. Also read from ENV.ALKSTFP_SECRET_ACCESS_KEY, ENV.AWS_SECRET_ACCESS_KEY and ENV.ALKS_SECRET_ACCESS_KEY, in that order.
* `token` - (Optional) The session token from a valid STS session. Also read from ENV.ALKSTFP_SESSION_TOKEN, ENV.AWS_SESSION_TOKEN and ENV.ALKS_SESSION_TOKEN, in that order.
* `shared_credentials_file` - (Optional) The the path to the shared credentials file. Also read from ENV.ALKSTFP_SHARED_CREDENTIALS_FILE, ENV.AWS_SHARED_CREDENTIALS_FILE and ENV.ALKS_SHARED_CREDENTIALS_FILE, in that order.
* `profile` - (Optional) This is the AWS profile name as set in the shared credentials file. Also read from ENV.ALKSTFP_PROFILE, ENV.AWS_PROFILE and ENV.ALKS_PROFILE, in that order.
* `proxy_mode` - (Optional) How ALKS requests choose a proxy: `auto` (default), `environment` or `none`. See [Proxies](#proxies). Applies to ALKS requests only, never to AWS calls. Added in 2.9.0.
* `proxy_url` - (Optional) Send every ALKS request through this proxy, for example `http://host.docker.internal:9000`. `NO_PROXY` does not apply to it. Cannot be combined with `proxy_mode = "none"`. The URL has no path, query or fragment; percent-encode `/`, `?`, `#` and `@` in a proxy password. Added in 2.9.0.
* `assume_role` - (Optional) This is the role information to assume before making calling ALKS. This feature works the same as the assume_role feature of the AWS Terraform Provider.
    * `role_arn` - (Required) The Role ARN to assume for calling the ALKS API.
    * `session_name` - (Optional) The session name to provide to AWS when creating STS credentials. Please see the AWS SDK documentation for more information.
    * `external_id` - (Optional) The external identifier to provide to AWS when creating STS credentials. Please see the AWS SDK documentation for more information.
    * `policy` - (Optional) This specifies additional policy restrictions to apply to the resulting STS credentials beyond any existing inline or managed policies. Please see the AWS SDK documentation for more information.
* `default_tags` - (Optional) This block can hold a block of tags to add to all roles created by this provider
    * `tags` - (Optional) Block of key value pairs to add to all roles
* `ignore_tags` - (Optional) Can contain a list of tag keys or key prefixes to exclude from `terraform plan` and `terraform apply`.  This is for tags added outside of the alks provider that are managed externally
    * `keys` - (Optional) List of keys to ignore
    * `key_prefixes` - (Optional) List of key prefixes to ignore. Any key starting with a string in this list will be ignored. 

---
### Environment variables reference

Every provider argument can also come from the environment. **A value in the provider block wins for that argument**, and then that argument's environment variables are not read. Otherwise the first variable that is set wins, left to right. This is decided argument by argument, not for the identity as a whole: for example, a `profile` in the provider block does not stop `AWS_ACCESS_KEY_ID` from being used, because access keys take precedence over a profile.

| Argument | Environment variables, in precedence order |
|---|---|
| `url` | `ALKSTFP_URL`, `ALKS_URL` |
| `access_key` | `ALKSTFP_ACCESS_KEY_ID`, `AWS_ACCESS_KEY_ID`, `ALKS_ACCESS_KEY_ID` |
| `secret_key` | `ALKSTFP_SECRET_ACCESS_KEY`, `AWS_SECRET_ACCESS_KEY`, `ALKS_SECRET_ACCESS_KEY` |
| `token` | `ALKSTFP_SESSION_TOKEN`, `AWS_SESSION_TOKEN`, `ALKS_SESSION_TOKEN` |
| `profile` | `ALKSTFP_PROFILE`, `AWS_PROFILE`, `ALKS_PROFILE` |
| `shared_credentials_file` | `ALKSTFP_SHARED_CREDENTIALS_FILE`, `AWS_SHARED_CREDENTIALS_FILE`, `ALKS_SHARED_CREDENTIALS_FILE` |
| `account` | `ALKSTFP_ACCOUNT`, `ALKS_ACCOUNT`, `Account` |
| `role` | `ALKSTFP_ROLE`, `ALKS_ROLE`, `Role` |
| `proxy_mode` | `ALKSTFP_PROXY_MODE`, `ALKS_PROXY_MODE` |
| `proxy_url` | `ALKSTFP_PROXY_URL`, `ALKS_PROXY_URL` |

**Conflicts are errors.** If two of the provider's own names for the same argument are set (`ALKSTFP_*`, `ALKS_*`, or the legacy `Account`/`Role`), the provider stops with an error naming them, even if their values are equal. Unset one, or set the argument in the provider block. `AWS_*` names are not part of this check, so `AWS_ACCESS_KEY_ID` and `ALKS_ACCESS_KEY_ID` together is fine; `AWS_*` wins.

**Newly read in 2.9.0:** `ALKS_PROFILE`, `ALKS_SHARED_CREDENTIALS_FILE`, `ALKS_ACCOUNT` and `ALKS_ROLE` were ignored before 2.9.0. If your shell already exports them for other ALKS tooling, the provider now uses them. `ALKS_ACCOUNT` with `ALKS_ROLE` switches it to that account and role. Unset them, or set the argument in the provider block, to keep the earlier result.

`Account` and `Role` are the original names and keep working. On macOS and Linux they are case-sensitive (`ACCOUNT` is a different variable).

**`ALKSTFP_*` gives this provider its own identity.**
- **Credentials are a set.** If any of `ALKSTFP_ACCESS_KEY_ID`, `ALKSTFP_SECRET_ACCESS_KEY` or `ALKSTFP_SESSION_TOKEN` is set, all three come from `ALKSTFP_`. A missing one stays empty; it is never filled from `AWS_*` or `ALKS_*`, so a key pair from one session is never combined with another session's token.
- **Not mixed with the provider block.** Setting `access_key`, `secret_key` or `token` in the provider block while an `ALKSTFP_` credential variable is set is an error. Put all the credentials in one place.
- **An incomplete set is an error.** If an `ALKSTFP_` credential is set but the access key or secret key is still missing, and no profile is set, the provider stops with an error. It does not fall back to `AWS_*` keys, `AWS_PROFILE` or the default profile. While an `ALKSTFP_` credential is set, the profile comes only from `ALKSTFP_PROFILE` (or the provider block), never from `AWS_PROFILE` or `ALKS_PROFILE`.
- **The profile wins.** If `ALKSTFP_PROFILE` is set, `AWS_*` and `ALKS_*` keys exported for other tooling are ignored, so the profile is used.
- **Scope.** These rules apply only when the `ALKSTFP_` variables are used. Otherwise the behaviour is as in earlier versions, where access keys, including `AWS_*`, take precedence over a profile, as in the AWS provider. A value in the provider block wins for its own argument.
- **`shared_credentials_file` is honoured throughout.** When a credentials file is given in the provider block or through `ALKSTFP_`/`ALKS_`, profiles are read only from that file, and a file that cannot be read is an error. Profiles are never silently read from the default `~/.aws/credentials`.

`access_key`, `secret_key`, `token` and `proxy_url` are marked sensitive, so Terraform does not show them.

### Proxies

From 2.9.0 the provider finds a proxy for its **ALKS** requests the way a browser does. It never changes how it calls AWS. With the default, `proxy_mode = "auto"`, each ALKS request uses:

1. `HTTPS_PROXY`, unless `NO_PROXY` matches the ALKS host. This is the same as earlier versions.
2. Otherwise the operating system's proxy settings: on **macOS** the system proxy and PAC file; on **Windows** the WinHTTP / Internet Options settings, PAC file and automatic detection. Linux has no OS step.
3. Otherwise a direct connection.

Rules:
- Only the first entry of a PAC answer is used. If that proxy does not answer, the request fails with an error naming it.
- `NO_PROXY` also exempts hosts from an OS proxy.
- `localhost` is never proxied.

This means Zscaler laptops that only route proxy-aware applications reach ALKS with nothing to configure.

| `proxy_mode` | Behaviour |
|---|---|
| `auto` (default) | As above |
| `environment` | `HTTPS_PROXY`/`NO_PROXY` only: exactly the behaviour before 2.9.0 |
| `none` | Always direct, ignoring every proxy setting |

`proxy_url` forces one proxy for every ALKS request.

**In a container** the host's proxy settings are not visible, so `auto` falls back to `HTTPS_PROXY`, or direct. With Docker Desktop, the Zscaler local proxy is reachable as `host.docker.internal`:

```hcl
provider "alks" {
  url       = "https://alks.example.com/rest"
  proxy_url = "http://host.docker.internal:9000"   # or set ALKS_PROXY_URL
}
```

`127.0.0.1:9000` does not work inside a container: it is the container itself. Setting `HTTPS_PROXY` instead also sends AWS calls through the proxy. `proxy_url`/`ALKS_PROXY_URL` affect ALKS only.

**Debugging:** with `TF_LOG=DEBUG`, the log has one line per client with the proxy decision, for example `alks-go: proxy mode=auto source=pac proxy=http://127.0.0.1:9000`. It never contains credentials.

### Supported Versions

| Terraform 0.10.x       |  Terraform 0.11.x        | Terraform 0.12.x - Terraform 0.15.x  | Terraform 1.0.x   |
| ---------------------- | ------------------------ | ------------------------------------ | ----------------- |
| ALKS TFP 0.9.0 < 1.3.0 | ALKS TFP 1.3.0 < 1.5.15  | ALKS TFP 1.3.0+                      | ALKS TFP 2.0.0+   |

For questions, please reach out to the [ALKS team](https://github.com/orgs/Cox-Automotive/teams/cai-internal-tools).