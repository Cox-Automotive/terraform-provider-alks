package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAlksLTKCreate(t *testing.T) {
	var resp alks.CreateIamUserResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAlksLtkDestroy(&resp),
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAlksLTKCreateConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("alks_ltk.foo", "iam_username", "TEST_LTK_USER")),
			},
		},
	})
}

func TestAlksLTKCreateWithTags(t *testing.T) {
	var resp alks.CreateIamUserResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAlksLtkDestroy(&resp),
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAlksLTKCreateWithTagsConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alks_ltk.foo", "iam_username", "TEST_LTK_USER"),
					resource.TestCheckResourceAttr("alks_ltk.foo", "tags.foo", "bar"),
					resource.TestCheckResourceAttr("alks_ltk.foo", "tags.cloud", "railway"),
				),
			},
		},
	})
}

func TestAlksLTKCreateWithTagsEmptyList(t *testing.T) {
	var resp alks.CreateIamUserResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAlksLtkDestroy(&resp),
		Steps: []resource.TestStep{
			// Create the resource
			{
				Config: testAlksLTKCreateWithTagsEmptyListConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("alks_ltk.foo", "iam_username", "TEST_LTK_USER"),
					resource.TestCheckResourceAttr("alks_ltk.foo", "tags_all.%", "0"),
				),
			},
		},
	})
}

func TestAccAlksLTKCreate_DefaultTags(t *testing.T) {
	var resp alks.CreateIamUserResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAlksLtkDestroy(&resp),
		Steps: []resource.TestStep{
			{
				// create resource with tags
				Config: testAccCheckAlksLtkCreateWithDefaultTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags.cloud", "railway"),
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags_all.cloud", "railway"),
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags_all.defaultTagKey1", "defaultTagValue1"),
				),
			},
		},
	})
}

func TestAccAlksLTKUpdate_DefaultTags(t *testing.T) {
	var resp alks.CreateIamUserResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAlksLtkDestroy(&resp),
		Steps: []resource.TestStep{
			{
				// create resource with tags
				Config: testAccCheckAlksLtkCreateWithDefaultTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags.cloud", "railway"),
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags_all.cloud", "railway"),
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags_all.defaultTagKey1", "defaultTagValue1"),
				),
			},
			{
				// update resource with new tags
				Config: testAccCheckAlksLtkUpdateWithDefaultTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags.cloud2", "railway2"),
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags_all.cloud2", "railway2"),
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags_all.%", "1"),
				),
			},
		},
	})
}

func TestAccAlksLTK_IgnoreTags(t *testing.T) {
	var resp alks.CreateIamUserResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAlksLtkDestroy(&resp),
		Steps: []resource.TestStep{
			{
				// create resource with tags
				Config: testAccCheckAlksLtkCreateWithDefaultTags,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags.cloud", "railway"),
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags_all.cloud", "railway"),
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags_all.defaultTagKey1", "defaultTagValue1"),
				),
			},
			{
				//Add tags externally.  These should not trigger an update because they are excluded by ignore_tags
				PreConfig: func() {
					client := testAccProvider.Meta().(*AlksClient).client
					tags := TagMap{
						"defaultTagKey1":        "defaultTagValue1",
						"cloud":                 "railway",
						"ignorePrefix:testKey1": "testValue1",
						"ignoreFullKey":         "testValue1",
					}
					userName := "TEST_LTK_USER"

					tagSlice := tagMapToSlice(tags)
					options := alks.IamUserOptions{
						IamUserName: &userName,
						Tags:        &tagSlice,
					}

					if _, err := client.UpdateIamUser(&options); err != nil {
						log.Printf("[INFO] Error in UpdateRole from test")
						return
					}
				},
				Config:   testAccCheckAlksLtkUpdateWithTagsWithIgnoredTags,
				PlanOnly: true, //This PlanOnly ensures there are no changes happening on this step.  Any changes will cause the test to error out because of uncompleted plan
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags_all.defaultTagKey1", "defaultTagValue1"),
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags.cloud", "railway"),
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags_all.cloud", "railway"),
					resource.TestCheckResourceAttr(
						"alks_ltk.foo", "tags_all.%", "2"),
				),
			},
		},
	})
}

func testAlksLtkDestroy(ltk *alks.CreateIamUserResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		providerStruct := testAccProvider.Meta().(*AlksClient)
		client := providerStruct.client

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alks_ltk" {
				continue
			}

			resp, err := client.GetIamUser(rs.Primary.ID)
			if resp != nil {
				return fmt.Errorf("Iam User still exists: %#v (%v)", resp, err)
			}
		}

		return nil
	}
}

const testAlksLTKCreateConfig = `
  resource "alks_ltk" "foo" {
    iam_username = "TEST_LTK_USER"
  }
`

const testAlksLTKCreateWithTagsConfig = `
  resource "alks_ltk" "foo" {
    iam_username = "TEST_LTK_USER"
	tags = {
		foo = "bar"
		cloud = "railway"
	  }
  }
`

const testAlksLTKCreateWithTagsEmptyListConfig = `
  resource "alks_ltk" "foo" {
    iam_username = "TEST_LTK_USER"
	tags = {}
  }
`

const testAccCheckAlksLtkCreateWithDefaultTags = `
	provider "alks" {
		default_tags {
			tags = {
				defaultTagKey1 = "defaultTagValue1"
			}
		}
	}
	resource "alks_ltk" "foo" {
		iam_username = "TEST_LTK_USER"
		tags = {
			cloud = "railway"
		  }
	  }
`

const testAccCheckAlksLtkUpdateWithDefaultTags = `
	provider "alks" {
	}
	resource "alks_ltk" "foo" {
		iam_username = "TEST_LTK_USER"
		tags = {
			cloud2 = "railway2"
		  }
	  }
`

const testAccCheckAlksLtkUpdateWithTagsWithIgnoredTags = `
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
	resource "alks_ltk" "foo" {
		iam_username = "TEST_LTK_USER"
		tags = {
			cloud = "railway"
		  }
	  }
`
