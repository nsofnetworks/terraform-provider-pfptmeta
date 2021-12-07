package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
