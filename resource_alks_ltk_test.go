package main

import (
	"fmt"
	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAlksLTKCreate(t *testing.T) {
	var resp alks.CreateLongTermKeyResponse

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
			// Update the resource
			{
				Config: testAlksLTKUpdateConfig,
				Check:  resource.ComposeTestCheckFunc(resource.TestCheckResourceAttr("alks_ltk.foo", "iam_username", "TEST_LTK_USER_2")),
						},
		},
	})
}

func testAlksLtkDestroy(ltk *alks.CreateLongTermKeyResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*alks.Client)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alks_ltk" {
				continue
			}

			resp, err := client.GetLongTermKey(rs.Primary.ID)
			if resp != nil {
				return fmt.Errorf("long term key still exists: %#v (%v)", resp, err)
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

const testAlksLTKUpdateConfig = `
  resource "alks_ltk" "foo" {
    iam_username = "TEST_LTK_USER_2"
  }
`
