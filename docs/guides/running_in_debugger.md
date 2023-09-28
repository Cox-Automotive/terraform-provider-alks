Steps for local TFP development:
1. Modify the code in the local terraform provider to run in debug mode.  Edit main.go:
```package main
import (
	"flag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	// For debugging provider, set the following environment variable
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	plugin.Serve(&plugin.ServeOpts{
		Debug: debug,
		ProviderAddr: "cox-automotive/alks",
		ProviderFunc: func() *schema.Provider {
			return Provider()
		},
	})
}
```
2.  Modify `~/.terraformrc` and configure dev_overrides to point to your code:
   ```provider_installation {
  dev_overrides {
    "cox-automotive/alks" = "/Users/james.barcelo/code/github.com/jcarlson/terraform-provider-alks"
   }

  direct {}
}

```
3.  Build local TFP in root dir (where the dev_override points):
   `cd ~/code/github.com/jcarlson/terraform-provider-alks && go build -o terraform-provider-alks -gcflags '-N -l'`
4. Setup launch.json config in vscode.  Note the args array uses your new debug flag:
```
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Debug ALKS TFP in local repo",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${workspaceFolder}",
            "env": {},
            "args": [
                "-debug"
            ],
            "envFile": "${workspaceFolder}/.vscode/private.env"
        }
    ]
}
```
5. Setup private.env as referenced in launch.json:
```
AWS_ACCESS_KEY_ID=ASIA3XTHISISPUBLICDATA
AWS_SECRET_ACCESS_KEY=KEEPITSECRET
AWS_SESSION_TOKEN=KEEPITSAFE
```
5. Run the launch config and then click on the "DEBUG CONSOLE".  You should see something like this:
```
TF_REATTACH_PROVIDERS='{"cox-automotive/alks":{"Protocol":"grpc","ProtocolVersion":5,"Pid":98763,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/jq/g4z2vsvn77s0c49s3zvk4cgw0000gq/T/plugin3567641228"}}}'
```
This means your custom provider is running and ready for connections.  Terraform connects with providers over grpc so the provider will run in server debug mode until you stop it.  You don't have to restart after every terraform run.
6. In seperate console run your terraform prefaced with the above environment variables:
```
TF_REATTACH_PROVIDERS='{"cox-automotive/alks":{"Protocol":"grpc","ProtocolVersion":5,"Pid":98763,"Test":true,"Addr":{"Network":"unix","String":"/var/folders/jq/g4z2vsvn77s0c49s3zvk4cgw0000gq/T/plugin3567641228"}}}' terraform apply
```
7. If you did everything right it is highly likely you will be able to successfully set breakpoints in vscode.
