package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const (
	metaportClusterAttachmentDependencies = `
resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service"
  mapped_service = "mapped.service.com"
}

resource "pfptmeta_network_element" "mapped-service2" {
  name           = "mapped service"
  mapped_service = "mapped.service2.com"
}

resource "pfptmeta_metaport_cluster" "metaport_cluster" {
  name = "metaport_cluster name1"
}
`
	metaportClusterAttachment1 = `
resource "pfptmeta_metaport_cluster_mapped_elements_attachment" "attachment" {
  metaport_cluster_id     = pfptmeta_metaport_cluster.metaport_cluster.id
  mapped_elements = [pfptmeta_network_element.mapped-service.id]
}
`
	metaportClusterAttachment2 = `
resource "pfptmeta_metaport_cluster_mapped_elements_attachment" "attachment2" {
  metaport_cluster_id     = pfptmeta_metaport_cluster.metaport_cluster.id
  mapped_elements = [pfptmeta_network_element.mapped-service2.id]
}
`
	metaport_clusterDataSource = `
data "pfptmeta_metaport_cluster" "metaport_cluster" {
  id = pfptmeta_metaport_cluster.metaport_cluster.id
}
`
)

func TestAccMetaportClusterAttachment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: metaportClusterAttachmentDependencies + metaportClusterAttachment1 + metaportClusterAttachment2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"pfptmeta_metaport_cluster_mapped_elements_attachment.attachment", "metaport_cluster_id",
						"pfptmeta_metaport_cluster.metaport_cluster", "id"),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_network_element.mapped-service", "id",
						"pfptmeta_metaport_cluster_mapped_elements_attachment.attachment", "mapped_elements.0"),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_network_element.mapped-service2", "id",
						"pfptmeta_metaport_cluster_mapped_elements_attachment.attachment2", "mapped_elements.0"),
				),
			},
			{
				Config: metaportClusterAttachmentDependencies + metaportClusterAttachment2,
			},
			{
				Config: metaportClusterAttachmentDependencies + metaportClusterAttachment2 + metaport_clusterDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.pfptmeta_metaport_cluster.metaport_cluster", "id",
						"pfptmeta_metaport_cluster.metaport_cluster", "id"),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_network_element.mapped-service2", "id",
						"data.pfptmeta_metaport_cluster.metaport_cluster", "mapped_elements.0"),
				),
			},
		},
	})
}
