package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMappedHost(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMappedHostStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_mapped_host.mapped-host", "network_element_id", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_mapped_host.mapped-host", "name", "host.test.com",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_mapped_host.mapped-host", "mapped_host", "host.test.com",
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_mapped_host.mapped-host-with-ipv4", "network_element_id", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_mapped_host.mapped-host-with-ipv4", "name", "host.test1.com",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_mapped_host.mapped-host-with-ipv4", "mapped_host", "10.72.1.0",
					),
				),
			},
		},
	})
}

func TestAccDataSourceMappedHost(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMappedHostStep1 + testAccDataSourceMappedHost,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_mapped_host.mapped-host", "network_element_id", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_mapped_host.mapped-host", "name", "host.test.com",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_mapped_host.mapped-host", "mapped_host", "host.test.com",
					),
				),
			},
		},
	})
}

const testAccResourceMappedHostStep1 = `
resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "mapped service name"
  description    = "some details about the mapped service"
  mapped_subnets = ["0.0.0.0/0", "10.20.30.0/24"]
}

resource "pfptmeta_mapped_host" "mapped-host" {
  network_element_id = pfptmeta_network_element.mapped-subnet.id
  mapped_host        = "host.test.com"
  name               = "host.test.com"
}

resource "pfptmeta_mapped_host" "mapped-host-with-ipv4" {
  network_element_id = pfptmeta_network_element.mapped-subnet.id
  mapped_host        = "10.72.1.0"
  name               = "host.test1.com"
}
`

const testAccDataSourceMappedHost = `

data "pfptmeta_mapped_host" "mapped-host" {
  network_element_id = pfptmeta_mapped_host.mapped-host.network_element_id
  name               = "host.test.com"
}`
