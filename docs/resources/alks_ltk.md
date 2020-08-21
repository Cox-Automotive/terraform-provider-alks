# Resource: alks_ltk

Creates a long term key user role in AWS.

## Example Usage

### ALKS IAM Role Creation
```hcl
resource "alks_ltk" "test_ltk_user" {
   iam_username             = "My_LTK_User_Name"
}
```

## Argument Reference

The following arguments are supported:
* `iam_username` - (Required) The name of the IAM user to create. This parameter allows a string of characters consisting of upper and lowercase alphanumeric characters with no spaces. You can also include any of the following characters: =,.@-. User names are not distinguished by case.
* `iam_user_arn` - (Computed) The ARN associated with the LTK user.
* `access_key` - (Computed) Generated access key for the LTK user. Note: This is saved in the state file, so please be aware of this.
* `secret_key` - (Computed) Generated secret key for the LTK user. Note: This is saved in the state file, so please be aware of this.

## Import

AWS IAM users can be imported using their user name, eg:
```
$ terraform import alks_ltk.my_ltk_resource MY_EXISTING_LTK_USER
```