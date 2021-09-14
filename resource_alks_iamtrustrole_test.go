package main

import (
	"testing"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlksIamTrustRole_Basic(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlksIamTrustRoleConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "name", "bar"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "type", "Inner Account"),
				),
			},
			{
				// update the resource
				Config: testAccCheckAlksIamTrustRoleConfigUpdateBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "name", "bar"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "type", "Inner Account"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "enable_alks_access", "true"),
				),
			},
		},
	})
}

const testAccCheckAlksIamTrustRoleConfigBasic = `
	resource "alks_iamrole" "foo" {
		name = "foo"
		type = "Amazon EC2"
		include_default_policies = false
	}

	resource "alks_iamtrustrole" "bar" {
		name = "bar"
		type = "Inner Account"
		trust_arn = "${alks_iamrole.foo.arn}"
	}
`

const testAccCheckAlksIamTrustRoleConfigUpdateBasic = `
	resource "alks_iamrole" "foo" {
		name = "foo"
		type = "Amazon EC2"
		include_default_policies = false
	}

	resource "alks_iamtrustrole" "bar" {
		name = "bar"
		type = "Inner Account"
		trust_arn = "${alks_iamrole.foo.arn}"
		enable_alks_access = true
	}
`
