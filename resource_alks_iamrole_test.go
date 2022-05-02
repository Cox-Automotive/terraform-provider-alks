package main

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlksIamRole_Basic(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlksIamRoleConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
				),
			},

			{
				// update the resource
				Config: testAccCheckAlksIamRoleConfigUpdateBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "enable_alks_access", "true"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "max_session_duration_in_seconds", "3600"),
				),
			},
		},
	})
}
func TestAccAlksIamRole_Tags(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				// create resource with tags
				Config: testAccCheckAlksIamRoleCreateWithTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "Car420"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags.testKey1", "testValue1"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags.testKey2", "testValue2"),
				),
			},
			{
				// update resource with tags
				Config: testAccCheckAlksIamRoleUpdateWithTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "Car420"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags.testKey3", "testValue3"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags.testKey2", "testValue2"),
				),
			},
		},
	})
}

func TestAccAlksIamRole_NoMaxDuration(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlksIamRoleConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
				),
			},
			{
				// update the resource
				Config: testAccCheckAlksIamRoleConfigUpdateNoMaxDuration,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "enable_alks_access", "true"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "max_session_duration_in_seconds", "3600"),
				),
			},
		},
	})
}

func TestAccIAMRole_NamePrefix(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlksIamRoleConfigNamePrefix,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.nameprefix", "name_prefix", "alks_test_acc_"),
					resource.TestMatchResourceAttr(
						"alks_iamrole.nameprefix", "name", regexp.MustCompile("alks_test_acc_[0-9]{26}")),
					resource.TestCheckResourceAttr(
						"alks_iamrole.nameprefix", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.nameprefix", "include_default_policies", "false"),
				),
			},
		},
	})
}

func TestAccIAMRole_NameAndNamePrefixConflict(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckAlksIamRoleConfigNameAndNamePrefixConflict,
				ExpectError: regexp.MustCompile(".*\"name\": conflicts with name_prefix.*"),
			},
		},
	})
}

func testAccCheckAlksIamRoleDestroy(role *alks.IamRoleResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*alks.Client)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alks_iamrole" && rs.Type != "alks_iamtrustrole" {
				continue
			}

			respz, err := client.GetIamRole(rs.Primary.ID)
			if respz != nil {
				return fmt.Errorf("Role still exists: %#v (%v)", respz, err)
			}
		}

		return nil
	}
}

func testAccCheckAlksIamRoleAttributes(role *alks.IamRoleResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[INFO] its this %v", role)
		if role.RoleName != "bar430" {
			return fmt.Errorf("Bad name: %s", role.RoleName)
		}
		if role.RoleType != "Amazon EC2" {
			return fmt.Errorf("Bad type: %s", role.RoleType)
		}

		return nil
	}
}

const testAccCheckAlksIamRoleConfigBasic = `
  resource "alks_iamrole" "foo" {
    name = "bar430"
    type = "Amazon EC2"
		include_default_policies = false
	}
`

const testAccCheckAlksIamRoleConfigUpdateBasic = `
	resource "alks_iamrole" "foo" {
		name = "bar430"
		type = "Amazon EC2"
		include_default_policies = false
		enable_alks_access = true
		max_session_duration_in_seconds = 3600
	}
`
const testAccCheckAlksIamRoleCreateWithTags = `
resource "alks_iamrole" "foo" {
  name = "bar430"
  type = "Amazon EC2"
	  include_default_policies = false
  }
  tags = {
	  testKey1 = "testValue1"
	  testKey2 = "testValue2"
  }
`

const testAccCheckAlksIamRoleUpdateWithTags = `
resource "alks_iamrole" "foo" {
  name = "bar430"
  type = "Amazon EC2"
	  include_default_policies = false
  }
  tags = {
	  testKey3 = "testValue3"
	  testKey2 = "testValue2"
  }
`

const testAccCheckAlksIamRoleConfigUpdateNoMaxDuration = `
	resource "alks_iamrole" "foo" {
		name = "bar430"
		type = "Amazon EC2"
		include_default_policies = false
		enable_alks_access = true
		max_session_duration_in_seconds = 3600
	}
`
const testAccCheckAlksIamRoleConfigNamePrefix = `
  resource "alks_iamrole" "nameprefix" {
    name_prefix = "alks_test_acc_"
    type = "Amazon EC2"
		include_default_policies = false
	}
`
const testAccCheckAlksIamRoleConfigNameAndNamePrefixConflict = `
  resource "alks_iamrole" "nameandnameprefixconflict" {
    name = "test-role"
    name_prefix = "alks_test_acc_"
    type = "Amazon EC2"
		include_default_policies = false
	}
`
