package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	routingGroupAttachmentDependencies = `
resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service"
  mapped_service = "mapped.service.com"
}

resource "pfptmeta_network_element" "mapped-service2" {
  name           = "mapped service"
  mapped_service = "mapped.service2.com"
}

resource "pfptmeta_routing_group" "routing_group" {
  name        = "routing group name"
}
`
	routingGroupAttachment1 = `
resource "pfptmeta_routing_group_mapped_elements_attachment" "attachment" {
  routing_group_id    = pfptmeta_routing_group.routing_group.id
  mapped_elements_ids = [pfptmeta_network_element.mapped-service.id]
}
`
	routingGroupAttachment2 = `
resource "pfptmeta_routing_group_mapped_elements_attachment" "attachment2" {
  routing_group_id    = pfptmeta_routing_group.routing_group.id
  mapped_elements_ids = [pfptmeta_network_element.mapped-service2.id]
}
`
)

func TestAccRoutingGroupAttachment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("routing_group", "v1/routing_groups"),
		Steps: []resource.TestStep{
			{
				Config: routingGroupAttachmentDependencies + routingGroupAttachment1 + routingGroupAttachment2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_routing_group_mapped_elements_attachment.attachment", "routing_group_id", regexp.MustCompile("^rg-.+$"),
					),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_network_element.mapped-service", "id",
						"pfptmeta_routing_group_mapped_elements_attachment.attachment", "mapped_elements_ids.0"),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_network_element.mapped-service2", "id",
						"pfptmeta_routing_group_mapped_elements_attachment.attachment2", "mapped_elements_ids.0"),
				),
			},
			{
				Config: routingGroupAttachmentDependencies + routingGroupAttachment2,
			},
			{
				Config: routingGroupAttachmentDependencies + routingGroupAttachment2 + routingGroupDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_routing_group.routing_group", "id", regexp.MustCompile("^rg-.+$"),
					),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_network_element.mapped-service2", "id",
						"data.pfptmeta_routing_group.routing_group", "mapped_elements_ids.0"),
				),
			},
		},
	})
}
