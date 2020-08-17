# Resource: alks_iamrole

Creates an custom ALKS IAM role for usage in an AWS account.

## Example Usage

### ALKS IAM Role Creation
```hcl
resource "alks_iamrole" "test_role" {
    name                     = "My_Test_Role"
    type                     = "Amazon EC2"
    include_default_policies = false
    enable_alks_access       = false
}
```

## Argument Reference

The following arguments are supported:
* `name` - (Required) The name of the ALKS IAM role which will be reflected in AWS and the ALKS UI.
* `type` - (Required) The role type to use. To see a list of available roles, [call this endpoint](https://pages.ghe.coxautoinc.com/ETS-CloudAutomation/ALKS-Documentation/#/aws-role-type-rest-service/getAllAwsRoleTypesUsingGET).
* `include_default_policies` - (Required) Whether or not the default manages policies should be attached to the role.
* `role_added_to_ip` - (Computed) Indicates whether or not an instance profile role was created.
* `arn` - (Computed) Provides the ARN of the role that was created.
* `ip_arn` - (Computed) If `role_added_to_ip` was `true` this will provide the ARN of the instance profile role.
* `enable_alks_access` - (Optional) If `true`, allows ALKS calls to be made by instance profiles or Lambda functions making use of this role. Note: This enables **machine identity** capability.