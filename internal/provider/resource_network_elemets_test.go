package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMappedSubnet(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      checkNetworkElementDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMappedSubnet,
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
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "net_id", regexp.MustCompile("^[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "tags.tag_name1", "tag_value1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "tags.tag_name2", "tag_value2",
					),
				),
			},
		},
	})
}

func TestAccResourceMappedService(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      checkNetworkElementDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceMappedService,
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
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-service", "net_id", regexp.MustCompile("^[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-service", "tags.tag_name1", "tag_value1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-service", "tags.tag_name2", "tag_value2",
					),
				),
			},
		},
	})
}

func checkNetworkElementDestroyed(s *terraform.State) error {
	c := provider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pfptmeta_network_element" {
			continue
		}
		neId := rs.Primary.ID
		_, err := client.GetNetworkElement(c, neId)
		if err == nil {
			return fmt.Errorf("network element %s still exists", neId)
		}
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("failed to verify network element %s was destroyed: %s", neId, err)
	}

	return nil
}

const testAccResourceMappedSubnet = `
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

const testAccResourceMappedService = `
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
