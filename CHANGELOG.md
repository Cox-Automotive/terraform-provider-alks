## 2.8.4 (July 9, 2026)

NOTES:

* Repository migrated from `github.com/Cox-Automotive/terraform-provider-alks` to `ghe.coxautoinc.com/ETS-CloudAutomation/terraform-provider-alks`. Development and releases continue from GHE; artifacts are still published to `github.com/Cox-Automotive/terraform-provider-alks` and the HashiCorp registry so consumer `source` addresses are unchanged.
* CI/CD workflows migrated from Travis CI to GitHub Actions on CAI self-hosted runners (`cai-standard-x86-linux`).
* Release workflow rewritten: GHE is now the source of truth; each release tags GHE, mirrors master and the new tag to github.com, then publishes via goreleaser to the HashiCorp registry.
* GPG signing key rotated. The original key (set up 2020, inaccessible since author's departure) was replaced with a new RSA 4096 key registered to the Cox-Automotive namespace on `registry.terraform.io`. Existing installs will see a one-time signature verification update on the next `terraform init` after upgrading.

## 2.8.2 (February 6, 2025)

* Removed unused local testing setup and associated `.gitignore` entries
* Switched goreleaser from deprecated `--rm-dist` to `--clean`
* Automated CI/CD pipeline with GitHub Actions (US852217)
