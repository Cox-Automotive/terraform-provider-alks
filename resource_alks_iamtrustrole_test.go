package main

import (
	"regexp"
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

func TestAccAlksIamTrustRole_NamePrefix(t *testing.T) {
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
						"alks_iamrole.nameprefix", "name_prefix", "alks_test_acc_"),
					resource.TestMatchResourceAttr(
						"alks_iamrole.nameprefix", "name", regexp.MustCompile("alks_test_acc_[0-9]{26}")),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "type", "Inner Account"),
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
const testAccCheckAlksIamTrustRoleConfigNamePrefix = `
	resource "alks_iamrole" "nameprefix_role" {
		name_prefix = "alks_test_acc_"
		type = "Amazon EC2"
		include_default_policies = false
	}

	resource "alks_iamtrustrole" "name-prefix_trustrole" {
		name_prefix = "alks_test_acc_"
		type = "Inner Account"
		trust_arn = "${alks_iamrole.nameprefix_role.arn}"
	}
`
