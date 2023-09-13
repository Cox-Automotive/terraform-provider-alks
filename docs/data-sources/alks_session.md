# Data Source: alks_session

Surfaces the AWS STS Credentials currently in use by the ALKS provider. Unlike the `alks_keys` data source, this data source does not
generate new credentials on-demand and instead outputs the credentials that the provider is already using.

## Example Usage

```hcl
data "alks_session" "current" {
   providers: alks.my_alias
}
```

## Argument Reference

* Note: This does not take any arguments. See below.

## Attribute Reference

* `access_key` - Current access key for the specified provider. If multiple providers, it takes the `provider` field. Otherwise, uses the initial provider.
* `secret_key` - Current secret key for the specified provider. If multiple providers, it takes the `provider` field. Otherwise, uses the initial provider.
* `session_token` - Current session token for the specified provider. If multiple providers, it takes the `provider` field. Otherwise, uses the initial provider.


## How it works 
- Whatever your default provider credentials are, will be used. If multiple providers have been configured, then one can point the data source to return keys for specific providers using `providers` field with an explicit alias.
