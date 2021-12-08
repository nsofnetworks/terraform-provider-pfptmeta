package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceRoutingGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("routing_group", "v1/routing_groups"),
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingGroupStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_routing_group.routing_group", "id", regexp.MustCompile("^rg-.+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_routing_group.routing_group", "name", "routing group name",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_routing_group.routing_group", "description", "routing group description",
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_routing_group.routing_group", "mapped_elements_ids.0", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_routing_group.routing_group", "sources.0", regexp.MustCompile("^grp-.+$"),
					),
				),
			},
			{
				Config: testAccRoutingGroupStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pfptmeta_routing_group.routing_group", "name", "routing group name1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_routing_group.routing_group", "description", "routing group description1",
					),
					resource.TestCheckNoResourceAttr("pfptmeta_routing_group.routing_group", "mapped_elements_ids"),
					resource.TestCheckNoResourceAttr("pfptmeta_routing_group.routing_group", "sources"),
				),
			},
		},
	})
}

func TestAccDataSourceRoutingGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("routing_group", "v1/routing_groups"),
		Steps: []resource.TestStep{
			{
				Config: testAccRoutingGroupStep1 + routingGroupDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_routing_group.routing_group", "id", regexp.MustCompile("^rg-.+$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_routing_group.routing_group", "name", "routing group name",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_routing_group.routing_group", "description", "routing group description",
					),
					resource.TestMatchResourceAttr(
						"data.pfptmeta_routing_group.routing_group", "mapped_elements_ids.0", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestMatchResourceAttr(
						"data.pfptmeta_routing_group.routing_group", "sources.0", regexp.MustCompile("^grp-.+$"),
					),
				),
			},
			{
				Config: testAccRoutingGroupStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pfptmeta_routing_group.routing_group", "name", "routing group name1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_routing_group.routing_group", "description", "routing group description1",
					),
					resource.TestCheckNoResourceAttr("pfptmeta_routing_group.routing_group", "mapped_elements_ids"),
					resource.TestCheckNoResourceAttr("pfptmeta_routing_group.routing_group", "sources"),
				),
			},
		},
	})
}

const testAccRoutingGroupStep1 = `
resource "pfptmeta_group" "group" {
  name        = "group name"
}

resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service"
  mapped_service = "mapped.service.com"
}

resource "pfptmeta_routing_group" "routing_group" {
  name = "routing group name"
  description = "routing group description"
  mapped_elements_ids = [pfptmeta_network_element.mapped-service.id]
  sources = [pfptmeta_group.group.id]
}
`

const testAccRoutingGroupStep2 = `
resource "pfptmeta_group" "group" {
  name        = "group name"
}

resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service"
  mapped_service = "mapped.service.com"
}

resource "pfptmeta_routing_group" "routing_group" {
  name = "routing group name1"
  description = "routing group description1"
}
`
const routingGroupDataSource = `
data "pfptmeta_routing_group" "routing_group" {
	id = pfptmeta_routing_group.routing_group.id
}
`
