package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMetaport(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("metaport", "v1/metaports"),
		Steps: []resource.TestStep{
			{
				Config: testAccMetaportStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_metaport.metaport", "id", regexp.MustCompile("^mp-[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport.metaport", "name", "metaport name",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport.metaport", "description", "metaport description",
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_metaport.metaport", "mapped_elements.0", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport.metaport", "allow_support", "true",
					),
				),
			},
			{
				Config: testAccMetaportStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport.metaport", "name", "metaport name1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport.metaport", "description", "metaport description1",
					),
					resource.TestCheckNoResourceAttr("pfptmeta_metaport.metaport", "mapped_elements"),
					resource.TestCheckResourceAttr("pfptmeta_metaport.metaport", "allow_support", "false"),
				),
			},
		},
	})
}

const testAccMetaportStep1 = `
resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "ms"
  mapped_subnets = ["0.0.0.0/0"]
}

resource "pfptmeta_metaport" "metaport" {
  name                  = "metaport name"
  description           = "metaport description"
  mapped_elements       = [pfptmeta_network_element.mapped-subnet.id]
}
`

const testAccMetaportStep2 = `
resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "ms"
  mapped_subnets = ["0.0.0.0/0"]
}

resource "pfptmeta_metaport" "metaport" {
  name                  = "metaport name1"
  description           = "metaport description1"
  allow_support         = false
}
`
