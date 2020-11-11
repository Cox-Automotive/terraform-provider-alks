---
page_title: "Example usage of ALKS TFP"
---

## Example

See [this example](https://github.com/Cox-Automotive/terraform-provider-alks/blob/master/examples/alks.tf) for a basic Terraform script which:

1. Creates an AWS provider and ALKS provider
   - Note: There are two ALKS / AWS providers to showcase multi-provider configuration in use.
2. Creates an IAM role via the ALKS provider
3. Attaches a policy to the created role using the AWS provider
4. Creates an LTK user via the ALKS provider.

This example is intended to show how to combine a typical AWS Terraform script with the ALKS provider to automate the creation of IAM roles and other infrastructure.
