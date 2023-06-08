package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMappedSubnet(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("network_element", "v1/network_elements"),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMappedSubnetStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "id", regexp.MustCompile("^ne-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "name", "mapped subnet name",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "description", "some details about the mapped subnet",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_subnets.0", "0.0.0.0/0",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_subnets.1", "10.20.30.0/24",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "type", "Mapped Subnet",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "tags.tag_name1", "tag_value1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "tags.tag_name2", "tag_value2",
					),
				),
			},
			{
				Config: testAccResourceMappedSubnetStep2,
			},
		},
	})
}

func TestAccResourceMappedService(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("network_element", "v1/network_elements"),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMappedServiceStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-service", "id", regexp.MustCompile("^ne-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-service", "name", "mapped service name",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-service", "description", "some details about the mapped service",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-service", "mapped_service", "mapped.service.com",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-service", "type", "Mapped Service",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-service", "tags.tag_name1", "tag_value1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-service", "tags.tag_name2", "tag_value2",
					),
				),
			},
			{
				Config: testAccResourceMappedServiceStep2,
			},
		},
	})
}

func TestAccDataSourceNetworkElement(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("network_element", "v1/network_elements"),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMappedSubnetStep1 + testAccDataSourceNetworkElement,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_network_element.mapped-subnet", "id", regexp.MustCompile("^ne-.*$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_network_element.mapped-subnet", "name", "mapped subnet name",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_network_element.mapped-subnet", "description", "some details about the mapped subnet",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_network_element.mapped-subnet", "mapped_subnets.0", "0.0.0.0/0",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_network_element.mapped-subnet", "mapped_subnets.1", "10.20.30.0/24",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_network_element.mapped-subnet", "type", "Mapped Subnet",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_network_element.mapped-subnet", "tags.tag_name1", "tag_value1",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_network_element.mapped-subnet", "tags.tag_name2", "tag_value2",
					),
				),
			},
		},
	})
}

const testAccResourceMappedSubnetStep1 = `
resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "mapped subnet name"
  description    = "some details about the mapped subnet"
  mapped_subnets = ["0.0.0.0/0", "10.20.30.0/24"]
  tags = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
}
`

const testAccResourceMappedSubnetStep2 = `
resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "mapped subnet name1"
  description    = "some details about the mapped subnet"
  mapped_subnets = ["0.0.0.0/0", "10.20.30.0/24"]
  tags = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
}
`

const testAccResourceMappedServiceStep1 = `
resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service name"
  description    = "some details about the mapped service"
  mapped_service = "mapped.service.com"
  tags = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
}
`

const testAccResourceMappedServiceStep2 = `
resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service name1"
  description    = "some details about the mapped service"
  mapped_service = "mapped.service.com"
  tags = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
}
`

const testAccDataSourceNetworkElement = `

data "pfptmeta_network_element" "mapped-subnet" {
  id = pfptmeta_network_element.mapped-subnet.id
}`
