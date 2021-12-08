package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccResourceRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("role", "v1/roles"),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRoleStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_role.admin_role", "id", regexp.MustCompile("^rol-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "name", "admin role",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "description", "role with all privileges",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "all_read_privileges", "true",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "all_write_privileges", "true",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "all_suborgs", "false",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.with_privileges", "name", "with privs",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.with_privileges", "privileges.0", "metaports:read",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.with_privileges", "privileges.1", "metaports:write",
					),
				),
			},
			{
				Config: testAccResourceRoleStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_role.admin_role", "id", regexp.MustCompile("^rol-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "name", "admin role1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "description", "role with all privileges1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "all_read_privileges", "false",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "all_write_privileges", "false",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "all_suborgs", "false",
					),
				),
			},
		},
	})
}

const testAccResourceRoleStep1 = `
resource "pfptmeta_role" "admin_role" {
  name                 = "admin role"
  description          = "role with all privileges"
  all_read_privileges  = true
  all_write_privileges = true
}

resource "pfptmeta_role" "with_privileges" {
  name                 = "with privs"
  privileges		   = ["metaports:read", "metaports:write"]
}
`

const testAccResourceRoleStep2 = `
resource "pfptmeta_role" "admin_role" {
  name                 = "admin role1"
  description          = "role with all privileges1"
  all_read_privileges  = false
}
`