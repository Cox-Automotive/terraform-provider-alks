package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	// "github.com/Cox-Automotive/alks-go"
)

func TestAccAlksSession_Basic(t *testing.T) {
	// var resp alks.SessionResponse

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckAlksSessionConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_session.foo", "name", "bar420"),
					resource.TestCheckResourceAttr(
						"alks_session.foo", "use_iam", "false"),
					resource.TestCheckResourceAttr(
						"alks_session.foo", "access_key", "asdf"),
					resource.TestCheckResourceAttr(
						"alks_session.foo", "secret_key", "qwerty"),
					resource.TestCheckResourceAttr(
						"alks_session.foo", "session_token", "zxcv"),
				),
			},
		},
	})
}

const testAccCheckAlksSessionConfig_basic = `
resource "alks_session" "foo" {
    use_iam = false
    name = "bar420"
    access_key = "asdf"
    secret_key = "qwerty"
    session_token = "zxcv"
}
`
