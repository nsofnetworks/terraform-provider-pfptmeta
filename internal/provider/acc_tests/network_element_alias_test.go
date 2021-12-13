package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceNetworkElementAlias(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAliasStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element_alias.alias", "network_element_id", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element_alias.alias", "alias", "alias.test.com",
					),
				),
			},
		},
	})
}

func TestAccDataSourceNetworkElementAlias(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAliasStep1 + testAccResourceAliasDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_network_element_alias.alias", "network_element_id", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_network_element_alias.alias", "alias", "alias.test.com",
					),
				),
			},
		},
	})
}

const testAccResourceAliasStep1 = `
resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service name"
  description    = "some details about the mapped service"
  mapped_service = "mapped.service.com"
}

resource "pfptmeta_network_element_alias" "alias" {
  network_element_id = pfptmeta_network_element.mapped-service.id
  alias = "alias.test.com"
}
`

const testAccResourceAliasDataSource = `

data "pfptmeta_network_element_alias" "alias" {
  network_element_id = pfptmeta_network_element_alias.alias.network_element_id
  alias              = "alias.test.com"
}
`
