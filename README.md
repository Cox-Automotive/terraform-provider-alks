# ALKS Provider for Terraform

# Building & Testing

```
make
```

# Sample STS

```
alks sessions open -i -a "######/ALKSAdmin - awsaepnp" -r "Admin" -o creds -n "tf-iam" -f
alks sessions open -a "######/ALKSPowerUser - awsaepnp" -r "PowerUser" -o creds -n "tf-main" -f
```


## Common build issues

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