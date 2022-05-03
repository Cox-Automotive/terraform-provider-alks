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

func TestAccAlksIamTrustRole_Tags(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlksIamTrustRoleCreateWithTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "name", "bar"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "type", "Inner Account"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "tags.testKey1", "testValue1"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "tags.testKey2", "testValue2"),
				),
			},
			{
				// update the resource
				Config: testAccCheckAlksIamTrustRoleUpdateWithTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "name", "bar"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "type", "Inner Account"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "tags.testKey3", "testValue3"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "tags.testKey2", "testValue2"),
				),
			},
		},
	})
}

func TestAccAlksIamTrustRole_DefaultTags(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlksIamRoleTrustCreateWithTagsWithDefault,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "name", "bar"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "type", "Inner Account"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "tags.testKey1", "testValue1"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "tags.testKey2", "testValue2"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "tags_all.defaultTagKey1", "defaultTagValue1"),
				),
			},
			{
				// update the resource
				Config: testAccCheckAlksIamTrustRoleUpdateWithTagsWithDefault,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "name", "bar"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "type", "Inner Account"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "tags.testKey3", "testValue3"),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.bar", "tags_all.defaultTagKey2", "defaultTagValue2"),
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
				Config: testAccCheckAlksIamTrustRoleConfigNamePrefix,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.nameprefix_trustrole", "name_prefix", "alks_test_acc_"),
					resource.TestMatchResourceAttr(
						"alks_iamtrustrole.nameprefix_trustrole", "name", regexp.MustCompile("alks_test_acc_[0-9]{26}")),
					resource.TestCheckResourceAttr(
						"alks_iamtrustrole.nameprefix_trustrole", "type", "Inner Account"),
				),
			},
		},
	})
}

func TestAccAlksIamTrustRole_NameAndNamePrefixConflict(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckAlksIamTrustRoleConfigNameAndNamePrefixConflict,
				ExpectError: regexp.MustCompile(".*\"name\": conflicts with name_prefix.*"),
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
const testAccCheckAlksIamTrustRoleCreateWithTags = `
	resource "alks_iamrole" "foo" {
		name = "foo"
		type = "Amazon EC2"
		include_default_policies = false
	}
	resource "alks_iamtrustrole" "bar" {
		name = "bar"
		type = "Inner Account"
		trust_arn = "${alks_iamrole.foo.arn}"
		tags = {
			testKey1 = "testValue1"
			testKey2 = "testValue2"
		}
	}
`

const testAccCheckAlksIamRoleTrustCreateWithTagsWithDefault = `
	provider "alks" {
		default_tags {
			tags = {
				defaultTagKey1 = "defaultTagValue1"
			}
		}
	}
	resource "alks_iamrole" "foo" {
		name = "foo"
		type = "Amazon EC2"
		include_default_policies = false

	}
	resource "alks_iamtrustrole" "bar" {
		name = "bar"
		type = "Inner Account"
		trust_arn = "${alks_iamrole.foo.arn}"
		tags = {
			testKey1 = "testValue1"
			testKey2 = "testValue2"
		}
	}

`
const testAccCheckAlksIamTrustRoleUpdateWithTagsWithDefault = `
	provider "alks" {
		default_tags {
			tags = {
				defaultTagKey2 = "defaultTagValue2"
			}
		}
	}
	resource "alks_iamrole" "foo" {
		name = "foo"
		type = "Amazon EC2"
		include_default_policies = false

	}
	resource "alks_iamtrustrole" "bar" {
		name = "bar"
		type = "Inner Account"
		trust_arn = "${alks_iamrole.foo.arn}"
		tags = {
			testKey1 = "testValue1"
			testKey3 = "testValue3"
		}
	}
`

const testAccCheckAlksIamTrustRoleUpdateWithTags = `
resource "alks_iamrole" "foo" {
	name = "foo"
	type = "Amazon EC2"
	include_default_policies = false
}
resource "alks_iamtrustrole" "bar" {
	name = "bar"
	type = "Inner Account"
	trust_arn = "${alks_iamrole.foo.arn}"
	tags = {
		testKey3 = "testValue3"
		testKey2 = "testValue2"
	}
}
`

const testAccCheckAlksIamTrustRoleConfigNamePrefix = `
	resource "alks_iamrole" "nameprefix_role" {
		name_prefix = "alks_test_acc_"
		type = "Amazon EC2"
		include_default_policies = false
	}

	resource "alks_iamtrustrole" "nameprefix_trustrole" {
		name_prefix = "alks_test_acc_"
		type = "Inner Account"
		trust_arn = "${alks_iamrole.nameprefix_role.arn}"
	}
`

const testAccCheckAlksIamTrustRoleConfigNameAndNamePrefixConflict = `
	resource "alks_iamrole" "nameprefixconflict_role" {
		name_prefix = "alks_test_acc_"
		type = "Amazon EC2"
		include_default_policies = false
	}

	resource "alks_iamtrustrole" "nameprefixconflict_trustrole" {
        name = "alks_test_acc"
		name_prefix = "alks_test_acc_"
		type = "Inner Account"
		trust_arn = "${alks_iamrole.nameprefixconflict_role.arn}"
	}
`
