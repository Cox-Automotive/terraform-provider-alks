# Resource: alks_iamtrustrole

Creates an IAM trust role on AWS.

## Example Usage


### ALKS IAM Role Creation
```hcl
resource "alks_iamtrustrole" "test_trust_role" {
    name                     = "My_Cross_Test_Role"
    type                     = "Cross Account"
    # type                   = "Inner Account"
    trust_arn                = "arn:aws:iam::123456789123:role/acct-managed/TestTrustRole"
    enable_alks_access       = false
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
* `name` - (Optional/Computed) The name of the IAM role to create. This parameter allows a string of characters consisting of upper and lowercase alphanumeric characters with no spaces. You can also include any of the following characters: =,.@-. Role names are not distinguished by case.
* `name_prefix` - (Optional/Computed) A prefix for a generated name of the IAM role to create. This parameter allows a string of characters consisting of upper and lowercase alphanumeric characters with no spaces. You can also include any of the following characters: =,.@-. Role names are not distinguished by case.
* `type` - (Required) 	The role type to use `Cross Account` or `Inner Account`.
* `trust_arn` - (Required) Account role ARN to trust.
  * _Note: This only allows **ONE** account role ARN. This is an intended security control by CAI._
* `role_added_to_ip` - (Computed) Indicates whether or not an instance profile role was created.
* `arn` - (Computed) Provides the ARN of the role that was created.
* `ip_arn` - (Computed) If `role_added_to_ip` was `true` this will provide the ARN of the instance profile role.
* `enable_alks_access` - (Optional) If `true`, allows ALKS calls to be made by instance profiles or Lambda functions making use of this role. Note: This enables **machine identity** capability.
* `tags` - (Optional) If present, will add specified tags onto role. 
* `max_session_duration_in_seconds` - (Optional) If present, will set maximum duration for role. Change forces re-creation of resource.


## Import

!> **Warning:** This will force-replace the resource.

ALKS IAM trust roles can be imported using the `name`, e.g.
```
$ terraform import alks_iamtrustrole.test_trust_role My_Cross_Test_Role
```