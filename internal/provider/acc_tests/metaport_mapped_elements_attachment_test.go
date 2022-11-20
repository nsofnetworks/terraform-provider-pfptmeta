package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	metaportAttachmentDependencies = `
resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service"
  mapped_service = "mapped.service.com"
}

resource "pfptmeta_network_element" "mapped-service2" {
  name           = "mapped service"
  mapped_service = "mapped.service2.com"
}

resource "pfptmeta_metaport" "metaport" {
  name = "metaport name1"
}
`
	metaportAttachment1 = `
resource "pfptmeta_metaport_mapped_elements_attachment" "attachment" {
  metaport_id     = pfptmeta_metaport.metaport.id
  mapped_elements = [pfptmeta_network_element.mapped-service.id]
}
`
	metaportAttachment2 = `
resource "pfptmeta_metaport_mapped_elements_attachment" "attachment2" {
  metaport_id     = pfptmeta_metaport.metaport.id
  mapped_elements = [pfptmeta_network_element.mapped-service2.id]
}
`
	metaportDataSource = `
data "pfptmeta_metaport" "metaport" {
  id = pfptmeta_metaport.metaport.id
}
`
)

func TestAccMetaportAttachment(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("metaport", "v1/metaports"),
		Steps: []resource.TestStep{
			{
				Config: metaportAttachmentDependencies + metaportAttachment1 + metaportAttachment2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_metaport_mapped_elements_attachment.attachment", "metaport_id", regexp.MustCompile("^mp-.+$"),
					),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_network_element.mapped-service", "id",
						"pfptmeta_metaport_mapped_elements_attachment.attachment", "mapped_elements.0"),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_network_element.mapped-service2", "id",
						"pfptmeta_metaport_mapped_elements_attachment.attachment2", "mapped_elements.0"),
				),
			},
			{
				Config: metaportAttachmentDependencies + metaportAttachment2,
			},
			{
				Config: metaportAttachmentDependencies + metaportAttachment2 + metaportDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_metaport.metaport", "id", regexp.MustCompile("^mp-.+$"),
					),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_network_element.mapped-service2", "id",
						"data.pfptmeta_metaport.metaport", "mapped_elements.0"),
				),
			},
		},
	})
}
