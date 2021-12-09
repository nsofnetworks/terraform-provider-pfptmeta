package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceRole,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.pfptmeta_role.admin_by_id", "name", "admin",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_role.admin_by_id", "all_suborgs", "true",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_role.admin_by_id", "all_write_privileges", "true",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_role.admin_by_id", "all_read_privileges", "true",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_role.admin_by_name", "id", "rol-Xqrzun95v8RA59E",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_role.admin_by_name", "all_suborgs", "true",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_role.admin_by_name", "all_write_privileges", "true",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_role.admin_by_name", "all_read_privileges", "true",
					),
				),
			},
		},
	})
}

const testAccDataSourceRole = `
data "pfptmeta_role" "admin_by_id" {
  id = "rol-Xqrzun95v8RA59E"
}

data "pfptmeta_role" "admin_by_name" {
  name = "admin"
}
`
