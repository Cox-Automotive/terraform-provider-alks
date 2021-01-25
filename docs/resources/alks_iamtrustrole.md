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

## Argument Reference

The following arguments are supported:
* `name` - (Required) The name of the IAM role to create. This parameter allows a string of characters consisting of upper and lowercase alphanumeric characters with no spaces. You can also include any of the following characters: =,.@-. Role names are not distinguished by case.
* `type` - (Required) 	The role type to use `Cross Account` or `Inner Account`.
* `trust_arn` - (Required) Account role ARN to trust.
* `role_added_to_ip` - (Computed) Indicates whether or not an instance profile role was created.
* `arn` - (Computed) Provides the ARN of the role that was created.
* `ip_arn` - (Computed) If `role_added_to_ip` was `true` this will provide the ARN of the instance profile role.
* `enable_alks_access` - (Optional) If `true`, allows ALKS calls to be made by instance profiles or Lambda functions making use of this role. Note: This enables **machine identity** capability.

## Import

ALKS IAM trust roles can be imported using the `name`, e.g.
```
$ terraform import alks_iamtrustrole.test_trust_role My_Cross_Test_Role
```