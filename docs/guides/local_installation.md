---
page_title: "Local Installation"
---

### Terraform Version < 0.13 Local Installation
* Download and install [Terraform](https://www.terraform.io/intro/getting-started/install.html)

* Download ALKS Provider binary for your platform from [Releases](https://github.com/Cox-Automotive/terraform-provider-alks/releases)

For example on macOS:

```
curl https://github.com/Cox-Automotive/terraform-provider-alks/releases/download/1.5.0/terraform-provider-alks_1.5.0_darwin_amd64.zip -O -J -L | unzip
```

* Configure Terraform to use this plugin by placing the binary in `.terraform.d/plugins/` on MacOS/Linux or `terraform.d\plugins\` in your user's "Application Data" directory on Windows.

* Note: If you've used a previous version of the ALKS provider and created a `.terraformrc` file in your home directory you'll want to remove it prior to updating.

### Terraform Version >= 0.13 Local Installation
* Download and install [Terraform](https://www.terraform.io/intro/getting-started/install.html)

* Download ALKS Provider binary for your platform from [Releases](https://github.com/Cox-Automotive/terraform-provider-alks/releases)
  
For example on macOS:

```
curl https://github.com/Cox-Automotive/terraform-provider-alks/releases/download/1.5.0/terraform-provider-alks_1.5.0_darwin_amd64.zip -O -J -L | unzip
```

* Go into the Terraform plugins path; `.terraform.d/plugins/` on MacOS/Linux or `terraform.d\plugins\` in your user's "Application Data" directory on Windows.

* Create the following directories: `coxautoinc.com/engineering-enablement/alks/1.5.0/<OS>_<ARCH>` and put the binary into the `<OS>_<ARCH>/` directory.
  * Note: This `<OS>_<ARCH>` will vary depending on your system. For example, 64-bit MacOS would be: `darwin_amd64` while 64-bit Windows 10 would be: `windows_amd64` 

* Finally, configure Terraform.
    * In your `versions.tf` or `main.tf` file you'll want to add the new ALKS provider as such:
    ```
    terraform {
        required_version = ">= 0.13"
        required_providers {
        alks = {
            source = "coxautoinc.com/engineering-enablement/alks"
            }
        }
    }
    ```

* Note: If you've previously installed our provider and it is stored in your remote state, you may need to run the [`replace-provider` command](https://www.terraform.io/docs/commands/state/replace-provider.html).
