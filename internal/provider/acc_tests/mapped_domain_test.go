package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMappedDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMappedDomainStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_mapped_domain.mapped-domain", "network_element_id", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_mapped_domain.mapped-domain", "name", "test.com",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_mapped_domain.mapped-domain", "mapped_domain", "test.com",
					),
				),
			},
		},
	})
}

func TestAccDataSourceMappedDomain(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMappedDomainStep1 + testAccDataSourceMappedDomain,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_mapped_domain.mapped-domain", "network_element_id", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_mapped_domain.mapped-domain", "name", "test.com",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_mapped_domain.mapped-domain", "mapped_domain", "test.com",
					),
				),
			},
		},
	})
}

const testAccResourceMappedDomainStep1 = `
resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "mapped service name"
  description    = "some details about the mapped service"
  mapped_subnets = ["0.0.0.0/0", "10.20.30.0/24"]
}

resource "pfptmeta_mapped_domain" "mapped-domain" {
  network_element_id = pfptmeta_network_element.mapped-subnet.id
  mapped_domain      = "test.com"
  name               = "test.com"
}
`

const testAccDataSourceMappedDomain = `

data "pfptmeta_mapped_domain" "mapped-domain" {
  network_element_id = pfptmeta_mapped_domain.mapped-domain.network_element_id
  name               = "test.com"
}`
