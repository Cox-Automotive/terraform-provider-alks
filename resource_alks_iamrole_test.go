package main

import (
    "fmt"
    "testing"

    "github.com/hashicorp/terraform/helper/resource"
    "github.com/hashicorp/terraform/terraform"
)

func TestAccAlksIamRole_Basic(t *testing.T) {
    var resp CreateRoleResponse

    resource.Test(t, resource.TestCase{
        PreCheck:     func() { testAccPreCheck(t) },
        Providers:    testAccProviders,
        CheckDestroy: testAccCheckAlksIamRoleDestroy,
        Steps: []resource.TestStep{
            resource.TestStep{
                Config: testAccCheckAlksIamRoleConfig_basic,
                Check: resource.ComposeTestCheckFunc(
                    testAccCheckAlksIamRoleExists("alks_iamrole.foo", &resp),
                    testAccCheckAlksIamRoleAttributes(&resp),
                    resource.TestCheckResourceAttr(
                        "alks_iamrole.foo", "name", "bar42"),
                    resource.TestCheckResourceAttr(
                        "alks_iamrole.foo", "type", "Amazon EC2"),
                    resource.TestCheckResourceAttr(
                        "alks_iamrole.foo", "include_default_policies", "false"),
                ),
            },
        },
    })
}

func testAccCheckAlksIamRoleDestroy(s *terraform.State) error {
    client := testAccProvider.Meta().(*AlksClient)

    for _, rs := range s.RootModule().Resources {
        if rs.Type != "alks_iamrole" {
            continue
        }

        resp, err := client.GetIamRole(rs.Primary.ID)

        if err == nil {
            return fmt.Errorf("Role still exists: %#v", resp)
        }
    }

    return nil
}


func testAccCheckAlksIamRoleExists(n string, role *CreateRoleResponse) resource.TestCheckFunc {
    return func(s *terraform.State) error {
        rs, ok := s.RootModule().Resources[n]

        if !ok {
            return fmt.Errorf("Not found: %s", n)
        }

        if rs.Primary.ID == "" {
            return fmt.Errorf("No role ID is set")
        }

        client := testAccProvider.Meta().(*AlksClient)

        foundRole, err := client.GetIamRoleByName(rs.Primary.Attributes["name"])

        if err != nil {
            return err
        }

        if foundRole.RoleArn != rs.Primary.ID {
            return fmt.Errorf("Role not found")
        }

        return nil
    }
}

func testAccCheckAlksIamRoleAttributes(role *CreateRoleResponse) resource.TestCheckFunc {
    return func(s *terraform.State) error {

        if role.RoleName != "bar42" {
            return fmt.Errorf("Bad name: %s", role.RoleName)
        }
        if role.RoleType != "Amazon EC2" {
            return fmt.Errorf("Bad type: %s", role.RoleType)
        }

        return nil
    }
}

const testAccCheckAlksIamRoleConfig_basic = `
resource "alks_iamrole" "foo" {
    name = "bar42"
    type = "Amazon EC2"
    include_default_policies = false
}
`