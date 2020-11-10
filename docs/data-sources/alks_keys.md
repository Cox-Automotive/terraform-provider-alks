# Data Source: alks_keys

Returns credentials for a given AWS account using ALKS.

## Example Usage

```hcl
data "alks_keys" "account_keys" {
   providers: alks.my_alias
}
```

## Argument Reference

* Note: This does not take any arguments. See below.

## Attribute Reference

* `access_key` - Generated access key for the specified provider. If multiple providers, it takes the `provider` field. Otherwise, uses the initial provider.
* `secret_key` - Generated secret key for the specified provider. If multiple providers, it takes the `provider` field. Otherwise, uses the initial provider.
* `session_token` - Generated session token for the specified provider. If multiple providers, it takes the `provider` field. Otherwise, uses the initial provider.
* `account` - The account number of the returned keys.
* `role` - The role from the returned keys.


## How it works 
- Whatever your default provider credentials are, will be used. If multiple providers have been configured, then one can point the data source to return keys for specific providers using `providers` field with an explicit alias.