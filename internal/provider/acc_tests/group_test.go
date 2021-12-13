package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccResourceGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("group", "v1/groups"),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_group.new_group", "id", regexp.MustCompile("^grp-.*$")),
					resource.TestCheckResourceAttr("pfptmeta_group.new_group", "name", "group name"),
					resource.TestCheckResourceAttr("pfptmeta_group.new_group", "description", "group description"),
					resource.TestCheckResourceAttr("pfptmeta_group.new_group", "expression", "tag_name:tag_value OR platform:macOS"),
				),
			},
			{
				Config: testAccResourceGroupStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_group.new_group", "id", regexp.MustCompile("^grp-.*$")),
					resource.TestCheckResourceAttr("pfptmeta_group.new_group", "name", "group name1"),
					resource.TestCheckResourceAttr("pfptmeta_group.new_group", "description", "group description1"),
					resource.TestCheckResourceAttr("pfptmeta_group.new_group", "expression", ""),
				),
			},
		},
	})
}

const testAccResourceGroupStep1 = `
resource "pfptmeta_group" "new_group" {
  name = "group name"
  description = "group description"
  expression = "tag_name:tag_value OR platform:macOS"
}
`

const testAccResourceGroupStep2 = `
resource "pfptmeta_group" "new_group" {
  name = "group name1"
  description = "group description1"
}
`

func TestAccDataSourceGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.pfptmeta_group.group_by_id", "id", "grp-Y5mxUn6ove1ybrn",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_group.group_by_id", "name", "test-group",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_group.group_by_id", "description", "Some group description",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_group.group_by_name", "id", "grp-Y5mxUn6ove1ybrn",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_group.group_by_name", "name", "test-group",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_group.group_by_name", "description", "Some group description",
					),
				),
			},
		},
	})
}

const testAccDataSourceGroup = `
data "pfptmeta_group" "group_by_id" {
  id = "grp-Y5mxUn6ove1ybrn"
}

data "pfptmeta_group" "group_by_name" {
  name = "test-group"
}
`
