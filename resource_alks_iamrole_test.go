package main

import (
	"fmt"
	"log"
	"regexp"
	"testing"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	awspolicy "github.com/hashicorp/awspolicyequivalence"
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
						"alks_iamrole.foo", "name", "bar430"),
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
						"alks_iamrole.foo", "name", "bar430"),
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

func TestAccAlksIamRole_DefaultTags_TrustPolicy(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				// create resource with tags
				Config: testAccCheckAlksIamRoleCreateWithTagsWithDefault_TrustPolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags.testKey1", "testValue1"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags.testKey2", "testValue2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags_all.defaultTagKey1", "defaultTagValue1"),
				),
			},
			{
				// update resource with tags
				Config: testAccCheckAlksIamRoleUpdateWithTagsWithDefault_TrustPolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags.testKey3", "testValue3"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags_all.defaultTagKey2", "defaultTagValue2"),
				),
			},
			{
				ResourceName:      "alks_iamrole.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAlksIamRole_DefaultTags_TrustPolicyUpdate(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				// create resource with tags
				Config: testAccCheckAlksIamRoleUpdateWithTagsWithDefault_TrustPolicy,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags_all.defaultTagKey2", "defaultTagValue2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags.testKey1", "testValue1"),
				),
			},
			{
				// update resource with tags
				Config: testAccCheckAlksIamRoleUpdateWithTagsWithDefault_TrustPolicyUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags.testKey4", "testValue4"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags_all.defaultTagKey1", "defaultTagValue1"),
					// Check the Assume Role Policy after the update
					testCheckAssumeRolePolicy("assume_role_policy", expectedAssumeRolePolicyAfterUpdate),
				),
			},
		},
	})
}

func testCheckAssumeRolePolicy(attr, expected string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		actual := s.RootModule().Resources["alks_iamrole.foo"].Primary.Attributes[attr]
		equivalent, err := awspolicy.PoliciesAreEquivalent(actual, expected)
		if err != nil {
			return fmt.Errorf("Unexpected error %s occured while comparing policies %s, and %s", err, expected, actual)
		}
		if !equivalent {
			return fmt.Errorf("Expected %s to be %s, got %s", attr, expected, actual)
		}
		return nil

	}
}

func TestAccAlksIamRole_DefaultTags_RoleType(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				// create resource with tags
				Config: testAccCheckAlksIamRoleCreateWithTagsWithDefault_RoleType,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags.testKey1", "testValue1"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags.testKey2", "testValue2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags_all.defaultTagKey1", "defaultTagValue1"),
				),
			},
			{
				// update resource with tags
				Config: testAccCheckAlksIamRoleUpdateWithTagsWithDefault_RoleType,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags.testKey3", "testValue3"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags_all.defaultTagKey2", "defaultTagValue2"),
				),
			},
		},
	})
}

func TestAccAlksIamRole_DefaultTagsEmpty(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				// create resource with tags
				Config: testAccCheckAlksIamRoleCreateWithTagsWithEmptyDefault,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
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
				Config: testAccCheckAlksIamRoleCreateWithTagsWithDefaultTagsEmpty,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
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
		},
	})
}

func TestAccAlksIamRole_IgnoreTagsEmpty(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				// create resource with tags
				Config: testAccCheckAlksIamRoleCreateWithTagsWithEmptyIgnore,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
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
				Config: testAccCheckAlksIamRoleUpdateWithTagsWithEmptyIgnore,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags_all.defaultTagKey1", "defaultTagValue1"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags.testKey2", "testValue2"),
				),
			},
		},
	})
}

func TestAccAlksIamRole_IgnoreTags(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				// create resource with tags
				Config: testAccCheckAlksIamRoleCreateWithTagsWithIgnoreTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags_all.defaultTagKey1", "defaultTagValue1"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags_all.testKey2", "testValue2"),
				),
			},
			{
				//Add tags externally.  These should not trigger an update because they are excluded by ignore_tags
				PreConfig: func() {
					client := testAccProvider.Meta().(*AlksClient).client
					tags := TagMap{
						"defaultTagKey1":        "defaultTagValue1",
						"testKey2":              "testValue2",
						"ignorePrefix:testKey1": "testValue1",
						"ignoreFullKey":         "testValue1",
					}
					roleName := "bar430"

					tagSlice := tagMapToSlice(tags)
					options := alks.UpdateIamRoleRequest{
						RoleName: &roleName,
						Tags:     &tagSlice,
					}
					if _, err := client.UpdateIamRole(&options); err != nil {
						fmt.Printf("Error in UpdateRole from test")
						return
					}
				},
				Config:   testAccCheckAlksIamRoleUpdateWithTagsWithIgnoredTags,
				PlanOnly: true, //This PlanOnly ensures there are no changes happening on this step.  Any changes will cause the test to error out because of uncompleted plan
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar430"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "tags_all.defaultTagKey1", "defaultTagValue1"),
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

func TestAccIAMRole_NameTooLong(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckAlksIamRoleConfigNameTooLong,
				ExpectError: regexp.MustCompile(".* expected length of name to be in the range \\(1 - 64\\).*"),
			},
		},
	})
}

func testAccCheckAlksIamRoleDestroy(role *alks.IamRoleResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerStruct := testAccProvider.Meta().(*AlksClient)
		client := providerStruct.client

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alks_iamrole" && rs.Type != "alks_iamtrustrole" {
				continue
			}

			respz, err := client.GetIamRole(rs.Primary.ID)
			if respz != nil {
				if respz.Exists == true {
					return fmt.Errorf("Role still exists: %#v (%v)", respz, err)
				}

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

func TestAccIAMRole_RoleTypeAndTrustPolicyBothPresent(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckAlksIamRoleBothRoleTypeAndTrustPolicyPresent,
				ExpectError: regexp.MustCompile(".*Error: Invalid combination.*"),
			},
		},
	})
}

func TestAccIAMRole_TrustPolicyAndIncludeDefaultPoliciesTrue(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				Config:      testAccCheckAlksIamRoleTrustPolicyAndDefaultPoliciesTrue,
				ExpectError: regexp.MustCompile(".*include_default_policies must be false or excluded if including an assume_role_policy.*"),
			},
		},
	})
}

func TestAccIAMRole_OnlyTrustPolicyPresent(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlksIamRoleWithOnlyTrustPolicyPresent,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.both_type_and_trust_policy", "name", "both_type_and_trust_policy"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.both_type_and_trust_policy", "include_default_policies", "false"),
				),
			},
			{
				ResourceName:      "alks_iamrole.both_type_and_trust_policy",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
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
		tags = {
			testKey1 = "testValue1"
			testKey2 = "testValue2"
		}
	}
`

const testAccCheckAlksIamRoleCreateWithTagsWithDefault_RoleType = `
	provider "alks" {
		default_tags {
			tags = {
				defaultTagKey1 = "defaultTagValue1"
			}
		}
	}
	resource "alks_iamrole" "foo" {
		name = "bar430"
		type = "Amazon EC2"
		include_default_policies = false
		tags = {
			testKey1 = "testValue1"
			testKey2 = "testValue2"
		}
	}
`

const testAccCheckAlksIamRoleCreateWithTagsWithDefault_TrustPolicy = `
	provider "alks" {
		default_tags {
			tags = {
				defaultTagKey1 = "defaultTagValue1"
			}
		}
	}
	resource "alks_iamrole" "foo" {

		name = "bar430"
		assume_role_policy       = jsonencode({
			Version = "2012-10-17",
			Statement = [
			{
				Action    = "sts:AssumeRole",
				Effect    = "Allow",
				Principal = {
				Service = "ec2.amazonaws.com"
				}
			}
			]
		})
		include_default_policies = false
		tags = {
			testKey1 = "testValue1"
			testKey2 = "testValue2"
		}
	}
`

const testAccCheckAlksIamRoleCreateWithTagsWithEmptyDefault = `
	provider "alks" {
		default_tags {
		}
	}
	resource "alks_iamrole" "foo" {
		name = "bar430"
		type = "Amazon EC2"
		include_default_policies = false
		tags = {
			testKey1 = "testValue1"
			testKey2 = "testValue2"
		}
	}
`
const testAccCheckAlksIamRoleCreateWithTagsWithDefaultTagsEmpty = `
	provider "alks" {
		default_tags {
			tags = {

			}
		}
	}
	resource "alks_iamrole" "foo" {
		name = "bar430"
		type = "Amazon EC2"
		include_default_policies = false
		tags = {
			testKey1 = "testValue1"
			testKey2 = "testValue2"
		}
	}
`

const testAccCheckAlksIamRoleCreateWithTagsWithEmptyIgnore = `
	provider "alks" {
		ignore_tags {
		}
	}
	resource "alks_iamrole" "foo" {
		name = "bar430"
		type = "Amazon EC2"
		include_default_policies = false
		tags = {
			testKey1 = "testValue1"
			testKey2 = "testValue2"
		}
	}
`
const testAccCheckAlksIamRoleUpdateWithTagsWithEmptyIgnore = `
	provider "alks" {
		default_tags {
			tags = {
				defaultTagKey1 = "defaultTagValue1"
			}
		}
		ignore_tags {
		}
	}
	resource "alks_iamrole" "foo" {
		name = "bar430"
		type = "Amazon EC2"
		include_default_policies = false
		tags = {
			testKey1 = "testValue1"
			testKey2 = "testValue2"
		}
	}
`

const testAccCheckAlksIamRoleCreateWithTagsWithIgnoreTags = `
	provider "alks" {
		default_tags {
			tags = {
				defaultTagKey1 = "defaultTagValue1"
			}
		}
	}
	resource "alks_iamrole" "foo" {
		name = "bar430"
		type = "Amazon EC2"
		include_default_policies = false
		tags = {
			testKey2 = "testValue2"
		}
	}
`

const testAccCheckAlksIamRoleUpdateWithTagsWithIgnoredTags = `
	provider "alks" {
		default_tags {
			tags = {
				defaultTagKey1 = "defaultTagValue1"
			}
		}
		ignore_tags {
			keys = ["ignoreFullKey"]
			key_prefixes = ["ignorePrefix"]
		}
	}
	resource "alks_iamrole" "foo" {
		name = "bar430"
		type = "Amazon EC2"
		include_default_policies = false
		tags = {
			testKey2 = "testValue2"
		}
	}
`
const testAccCheckAlksIamRoleUpdateWithTagsWithDefault_RoleType = `
	provider "alks" {
		default_tags {
			tags = {
				defaultTagKey2 = "defaultTagValue2"
			}
		}
	}
	resource "alks_iamrole" "foo" {
		name = "bar430"
		type = "Amazon EC2"
		include_default_policies = false
		tags = {
			testKey1 = "testValue1"
			testKey3 = "testValue3"
		}
	}
`

const testAccCheckAlksIamRoleUpdateWithTagsWithDefault_TrustPolicy = `
	provider "alks" {
		default_tags {
			tags = {
				defaultTagKey2 = "defaultTagValue2"
			}
		}
	}
	resource "alks_iamrole" "foo" {
		name = "bar430"
		assume_role_policy       = jsonencode({
			Version = "2012-10-17",
			Statement = [
			  {
				Action    = "sts:AssumeRole",
				Effect    = "Allow",
				Principal = {
				  Service = "ec2.amazonaws.com"
				}
			  }
			]
		  })
		include_default_policies = false
		tags = {
			testKey1 = "testValue1"
			testKey3 = "testValue3"
		}
	}
`

const testAccCheckAlksIamRoleUpdateWithTagsWithDefault_TrustPolicyUpdate = `
	provider "alks" {
		default_tags {
			tags = {
				defaultTagKey1 = "defaultTagValue1"
			}
		}
	}
	resource "alks_iamrole" "foo" {
		name = "bar430"
		assume_role_policy       = jsonencode({
			Version = "2012-10-17",
			Statement = [
			  {
				Action    = "sts:AssumeRole",
				Effect    = "Allow",
				Principal = {
				  Service = "lambda.amazonaws.com"
				}
			  }
			]
		  })
		include_default_policies = false
		tags = {
			testKey1 = "testValue1"
			testKey4 = "testValue4"
		}
	}
`

const expectedAssumeRolePolicyAfterUpdate = `{
	"Version": "2012-10-17",
	"Statement": [
		{
			"Action": "sts:AssumeRole",
			"Effect": "Allow",
			"Principal": {
				"Service": "lambda.amazonaws.com"
			}
		}
	]
}`

const testAccCheckAlksIamRoleUpdateWithTags = `
resource "alks_iamrole" "foo" {
	name = "bar430"
	type = "Amazon EC2"
	include_default_policies = false
	tags = {
		testKey3 = "testValue3"
		testKey2 = "testValue2"
	}
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

const testAccCheckAlksIamRoleConfigNameTooLong = `
  resource "alks_iamrole" "nametoolong" {
    name = "nameandnametoolongggggggggggggggggggggggggggggggggggggggggggggggg"
    type = "Amazon EC2"
		include_default_policies = false
	}
`

const testAccCheckAlksIamRoleBothRoleTypeAndTrustPolicyPresent = `
	resource "alks_iamrole" "both_type_and_trust_policy" {
		name                     = "both_type_and_trust_policy"
		include_default_policies = false
		type                     = "Amazon EC2"
		assume_role_policy             = jsonencode({
			Version = "2012-10-17",
			Statement = [
				{
					Action = "sts:AssumeRole",
					Effect = "Allow",
					Principal = {
						Service = "databrew.amazonaws.com"
					},
					Sid = ""
				}
			]
		})
	}
`

const testAccCheckAlksIamRoleWithOnlyTrustPolicyPresent = `
	resource "alks_iamrole" "both_type_and_trust_policy" {
		name                     = "both_type_and_trust_policy"
		include_default_policies = false
		assume_role_policy       = jsonencode({
			Version = "2012-10-17",
			Statement = [
				{
					Action = "sts:AssumeRole",
					Effect = "Allow",
					Principal = {
						Service = "databrew.amazonaws.com"
					},
					Sid = ""
				}
			]
		})
	}
`
const testAccCheckAlksIamRoleTrustPolicyAndDefaultPoliciesTrue = `
	resource "alks_iamrole" "both_type_and_trust_policy" {
		name                     = "both_type_and_trust_policy"
		include_default_policies = true
		assume_role_policy       = jsonencode({
			Version = "2012-10-17",
			Statement = [
				{
					Action = "sts:AssumeRole",
					Effect = "Allow",
					Principal = {
						Service = "databrew.amazonaws.com"
					},
					Sid = ""
				}
			]
		})
	}
`
