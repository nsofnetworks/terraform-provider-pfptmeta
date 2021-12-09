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
`
