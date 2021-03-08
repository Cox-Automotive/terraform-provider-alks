---
page_title: "Local Installation"
---

### Terraform Version 0.12 Local Installation
* Download and install [Terraform](https://www.terraform.io/intro/getting-started/install.html)

* Download ALKS Provider binary for your platform from [Releases](https://github.com/Cox-Automotive/terraform-provider-alks/releases)

One-liner download for macOS / Linux:
```
mkdir -p ~/.terraform.d/plugins &&
      curl -Ls https://api.github.com/repos/Cox-Automotive/terraform-provider-alks/releases/latest |
      jq -r ".assets[] | select(.browser_download_url | contains(\"$(uname -s | tr A-Z a-z)\")) | select(.browser_download_url | contains(\"amd64\")) | .browser_download_url" |
            xargs -n 1 curl -Lo ~/.terraform.d/plugins/terraform-provider-alks.zip &&
      pushd ~/.terraform.d/plugins/ &&
      unzip ~/.terraform.d/plugins/terraform-provider-alks.zip -d terraform-provider-alks-tmp &&
      mv terraform-provider-alks-tmp/terraform-provider-alks* . &&
      chmod +x terraform-provider-alks* &&
      rm -rf terraform-provider-alks-tmp &&
      rm -rf terraform-provider-alks.zip &&
      popd
```

* Configure Terraform to use this plugin by placing the binary in `.terraform.d/plugins/` on MacOS/Linux or `terraform.d\plugins\` in your user's "Application Data" directory on Windows.

* Note: If you've used a previous version of the ALKS provider and created a `.terraformrc` file in your home directory you'll want to remove it prior to updating.

### Terraform Version 0.13 Local Installation
* Download and install [Terraform](https://www.terraform.io/intro/getting-started/install.html)

* Download ALKS Provider binary for your platform from [Releases](https://github.com/Cox-Automotive/terraform-provider-alks/releases)
  
One-liner download for macOS / Linux:
```
mkdir -p ~/.terraform.d/plugins &&
      curl -Ls https://api.github.com/repos/Cox-Automotive/terraform-provider-alks/releases/latest |
      jq -r ".assets[] | select(.browser_download_url | contains(\"$(uname -s | tr A-Z a-z)\")) | select(.browser_download_url | contains(\"amd64\")) | .browser_download_url" |
            xargs -n 1 curl -Lo ~/.terraform.d/plugins/terraform-provider-alks.zip &&
      pushd ~/.terraform.d/plugins/ &&
      unzip ~/.terraform.d/plugins/terraform-provider-alks.zip -d terraform-provider-alks-tmp &&
      mv terraform-provider-alks-tmp/terraform-provider-alks* . &&
      chmod +x terraform-provider-alks* &&
      rm -rf terraform-provider-alks-tmp &&
      rm -rf terraform-provider-alks.zip &&
      popd
```
!> **Warning:** Your binary has been placed in `.terraform.d/plugins/`; this will **NOT** suffice for Terraform version 0.13. For more information, [read here](https://www.terraform.io/upgrade-guides/0-13.html#new-filesystem-layout-for-local-copies-of-providers).

* Go into the Terraform plugins path; `.terraform.d/plugins/` on MacOS/Linux or `terraform.d\plugins\` in your user's "Application Data" directory on Windows.

* Create the following directories: `coxautoinc.com/engineering-enablement/alks/<VERSION>/<OS>_<ARCH>` and put the binary into the `<OS>_<ARCH>/` directory.
  * Note: This `<OS>_<ARCH>` will vary depending on your system. For example, 64-bit MacOS will be: `darwin_amd64` while 64-bit Windows 10 will be: `windows_amd64` 

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

* Note: If you've previously installed our provider, and it is stored in your remote state: you may need to run the [`replace-provider` command](https://www.terraform.io/docs/commands/state/replace-provider.html).
