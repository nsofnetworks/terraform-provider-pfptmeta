package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	trustedNetworkResourceStep1 = `
resource "pfptmeta_trusted_network" "network" {
  name         = "trusted network name"
  description  = "trusted network description"
  apply_to_org = true
  criteria {
    external_ip_config {
      addresses_ranges = ["192.1.0.0/16"]
    }
  }
  criteria {
    resolved_address_config {
      addresses_ranges = ["192.1.0.0/16"]
      hostname         = "office.address11.com"
    }
  }
}
`
	trustedNetworkResourceStep2 = `
resource "pfptmeta_trusted_network" "network" {
  name         = "trusted network name1"
  description  = "trusted network description1"
  apply_to_org = true
  criteria {
    resolved_address_config {
      addresses_ranges = ["192.1.0.0/16"]
      hostname         = "office.address11.com"
    }
  }
  criteria {
    external_ip_config {
      addresses_ranges = ["192.1.0.0/16"]
    }
  }
}
`
	trustedNetworkDataSource = `
data "pfptmeta_trusted_network" "network" {
  id = pfptmeta_trusted_network.network.id
}
`
)

func TestAccResourceTrustedNetwork(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("trusted_network", "v1/trusted_networks"),
		Steps: []resource.TestStep{
			{
				Config: trustedNetworkResourceStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_trusted_network.network", "id", regexp.MustCompile("^tn-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network", "name", "trusted network name"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network", "description", "trusted network description"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network",
						"criteria.0.external_ip_config.0.addresses_ranges.0", "192.1.0.0/16"),
					resource.TestCheckNoResourceAttr("pfptmeta_trusted_network.network", "criteria.0.resolved_address_config"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network",
						"criteria.1.resolved_address_config.0.addresses_ranges.0", "192.1.0.0/16"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network",
						"criteria.1.resolved_address_config.0.hostname", "office.address11.com"),
					resource.TestCheckNoResourceAttr("pfptmeta_trusted_network.network", "criteria.1.external_ip_config"),
				),
			},
			{
				Config: trustedNetworkResourceStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_trusted_network.network", "id", regexp.MustCompile("^tn-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network", "name", "trusted network name1"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network", "description", "trusted network description1"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network",
						"criteria.1.external_ip_config.0.addresses_ranges.0", "192.1.0.0/16"),
					resource.TestCheckNoResourceAttr("pfptmeta_trusted_network.network", "criteria.1.resolved_address_config"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network",
						"criteria.0.resolved_address_config.0.addresses_ranges.0", "192.1.0.0/16"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network",
						"criteria.0.resolved_address_config.0.hostname", "office.address11.com"),
					resource.TestCheckNoResourceAttr("pfptmeta_trusted_network.network", "criteria.0.external_ip_config"),
				),
			},
		},
	})
}

func TestAccDataSourceTrustedNetwork(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("trusted_network", "v1/trusted_networks"),
		Steps: []resource.TestStep{
			{
				Config: trustedNetworkResourceStep1 + trustedNetworkDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_trusted_network.network", "id", regexp.MustCompile("^tn-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network", "name", "trusted network name"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network", "description", "trusted network description"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network",
						"criteria.0.external_ip_config.0.addresses_ranges.0", "192.1.0.0/16"),
					resource.TestCheckNoResourceAttr("pfptmeta_trusted_network.network", "criteria.0.resolved_address_config"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network",
						"criteria.1.resolved_address_config.0.addresses_ranges.0", "192.1.0.0/16"),
					resource.TestCheckResourceAttr("pfptmeta_trusted_network.network",
						"criteria.1.resolved_address_config.0.hostname", "office.address11.com"),
					resource.TestCheckNoResourceAttr("pfptmeta_trusted_network.network", "criteria.1.external_ip_config"),
				),
			},
		},
	})
}
