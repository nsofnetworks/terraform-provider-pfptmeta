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
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "net_id", regexp.MustCompile("^[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "tags.tag_name1", "tag_value1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "tags.tag_name2", "tag_value2",
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_domains.0.name", regexp.MustCompile("step1.test[\\d]*.com$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_domains.0.mapped_domain", regexp.MustCompile("step1.test[\\d]*.com$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_domains.0.name", regexp.MustCompile("step1.test[\\d]*.com$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_domains.0.mapped_domain", regexp.MustCompile("step1.test[\\d]*.com$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_hosts.0.name", regexp.MustCompile("step1.host[\\d]*.com$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_hosts.1.name", regexp.MustCompile("step1.host[\\d]*.com$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_hosts.0.mapped_host", "10.0.0.1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_hosts.1.mapped_host", "10.0.0.1",
					),
				),
			},
			{
				Config: testAccResourceMappedSubnetStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_domains.0.name", regexp.MustCompile("step2.test[\\d]*.com$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_domains.0.mapped_domain", regexp.MustCompile("step2.test[\\d]*.com$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_domains.0.name", regexp.MustCompile("step2.test[\\d]*.com$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_domains.0.mapped_domain", regexp.MustCompile("step2.test[\\d]*.com$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_hosts.0.name", regexp.MustCompile("step2.host[\\d]*.com$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-subnet", "mapped_hosts.1.name", regexp.MustCompile("step2.host[\\d]*.com$"),
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
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-service", "net_id", regexp.MustCompile("^[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-service", "tags.tag_name1", "tag_value1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.mapped-service", "tags.tag_name2", "tag_value2",
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-service", "aliases.0", regexp.MustCompile("^step1.mapped.service[12]+.com$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-service", "aliases.1", regexp.MustCompile("^step1.mapped.service[12]+.com$"),
					),
				),
			},
			{
				Config: testAccResourceMappedServiceStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-service", "aliases.0", regexp.MustCompile("^step2.mapped.service[12]+.com$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.mapped-service", "aliases.1", regexp.MustCompile("^step2.mapped.service[12]+.com$"),
					),
				),
			},
		},
	})
}

func TestAccResourceDevice(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      checkNetworkElementDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDevice,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.device", "id", regexp.MustCompile("^ne-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.device", "name", "acc tests device",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.device", "description", "some details about the device",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.device", "owner_id", "usr-LdjvfnK5713B8K1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.device", "platform", "Linux",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_network_element.device", "type", "Device",
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_network_element.device", "net_id", regexp.MustCompile("^[\\d]+$"),
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

const testAccResourceMappedSubnetStep1 = `
resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "mapped subnet name"
  description    = "some details about the mapped subnet"
  mapped_subnets = ["0.0.0.0/0", "10.20.30.0/24"]
  tags = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
  mapped_domains {
    name = "step1.test.com"
    mapped_domain = "step1.test.com"
  }
  mapped_domains {
    name          = "step1.test1.com"
    mapped_domain = "step1.test1.com"
  }
  mapped_hosts {
    name        = "step1.host.com"
    mapped_host = "10.0.0.1"
  }
  mapped_hosts {
    name        = "step1.host1.com"
    mapped_host = "10.0.0.1"
  }
}
`

const testAccResourceMappedSubnetStep2 = `
resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "mapped subnet name"
  description    = "some details about the mapped subnet"
  mapped_subnets = ["0.0.0.0/0", "10.20.30.0/24"]
  tags = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
  mapped_domains {
    name = "step2.test.com"
    mapped_domain = "step2.test.com"
  }
  mapped_domains {
    name = "step2.test1.com"
    mapped_domain = "step2.test1.com"
  }
  mapped_hosts {
    name        = "step2.host.com"
    mapped_host = "10.0.0.1"
  }
  mapped_hosts {
    name        = "step2.host1.com"
    mapped_host = "10.0.0.1"
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
  aliases = ["step1.mapped.service1.com", "step1.mapped.service2.com"]
}
`

const testAccResourceMappedServiceStep2 = `
resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service name"
  description    = "some details about the mapped service"
  mapped_service = "mapped.service.com"
  tags = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
  aliases = ["step2.mapped.service1.com", "step2.mapped.service2.com"]
}
`

const testAccResourceDevice = `
resource "pfptmeta_network_element" "device" {
  name        = "acc tests device"
  description = "some details about the device"
  owner_id    = "usr-LdjvfnK5713B8K1"
  platform    = "Linux"
  aliases     = ["alias.device.com"]
}`
