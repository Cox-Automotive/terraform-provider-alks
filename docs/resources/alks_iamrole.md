# Resource: alks_iamrole

Creates an custom ALKS IAM role for usage in an AWS account.

## Example Usage

### ALKS IAM Role Creation

#### IAM Role with a custom trust policy document

```hcl
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

This will create a role with the exact name `My_Test_Role`. Specifying a custom trust policy like this is currently only supported for single-service trust policies trusting an approved AWS service, and at the moment no extra fields may be provided such as the "Condition" or "Resource" keys. At this time, the only acceptable changes to the JSON string passed to the assume_role_policy field above are that `ec2.amazonaws.com` can be swapped out for any single approved service, and the `Sid` field may be omitted or populated with any valid Sid according to AWS's documentation.

#### IAM Role specifying a role type

```hcl
resource "alks_iamrole" "test_role" {
    name                     = "My_Test_Role"
    type                     = "Amazon EC2"
    include_default_policies = false
    enable_alks_access       = false
}
```

### ALKS IAM Role Creation with Name Prefix

```hcl
resource "alks_iamrole" "test_role" {
    name_prefix              = "My_Test_Role_"
    type                     = "Amazon EC2"
    include_default_policies = false
    enable_alks_access       = false
}
```

This will create a role named similar to `My_Test_Role_20211103145836382400000001`.
This is useful to avoid name conflicts when using the same terraform in multiple regions.

### ALKS Dynamic Role Creation

```hcl
resource "alks_iamrole" "test_dynamic_role" {
    name                     = "my_dynamic_role"
    type                     = "Amazon EKS IRSA"
    include_default_policies = false
    enable_alks_access       = false
    template_fields          = {
        "OIDC_PROVIDER"       = "abc1234"
        "K8S_NAMESPACE"       = "default"
        "K8S_SERVICE_ACCOUNT" = "my-service-account"
    }
}
```

### ALKS IAM Role Creation With Tags

```hcl
resource "alks_iamrole" "test_role" {
    name                     = "My_Test_Role"
    type                     = "Amazon EC2"
    include_default_policies = false
    enable_alks_access       = false
    tags                     = {
        "tagKey" = "tagValue"
    } 
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional/Computed) The name of the ALKS IAM role which will be reflected in AWS and the ALKS UI.
* `name_prefix` - (Optional/Computed) A prefix for a generated name of the ALKS IAM role which will be reflected in AWS and the ALKS UI.
* `type` - (Required) The role type to use. To see a list of available roles, [call this endpoint](https://pages.ghe.coxautoinc.com/ETS-CloudAutomation/ALKS-Documentation/#/aws-role-type-rest-service/getAllAwsRoleTypesUsingGET).
* `include_default_policies` - (Required) Whether or not the default manages policies should be attached to the role.
* `role_added_to_ip` - (Computed) Indicates whether or not an instance profile role was created.
* `arn` - (Computed) Provides the ARN of the role that was created.
* `ip_arn` - (Computed) If `role_added_to_ip` was `true` this will provide the ARN of the instance profile role.
* `enable_alks_access` - (Optional) If `true`, allows ALKS calls to be made by instance profiles or Lambda functions making use of this role. Note: This enables **machine identity** capability.
* `template_fields` - (Optional) If present, will submit template field data to ALKS.  Note: This will generate an error if the role type does not support template fields.
* `tags` - (Optional) If present, will add specified tags onto role. 

## Import

!> **Warning:** This will force-replace the resource.

ALKS IAM roles can be imported using the `name`, e.g.

```sh
terraform import alks_iamrole.test_role My_Test_Role
```