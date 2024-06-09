package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	ipNetworkStep1 = `
resource "pfptmeta_ip_network" "in" {
  name        = "in"
  description = "in desc"
  cidrs       = ["0.0.0.0/32"]
  countries   = ["US"]
}
`
	ipNetworkStep2 = `
resource "pfptmeta_ip_network" "in" {
  name        = "in 1"
  description = "in desc 1"
  cidrs       = ["192.5.0.0/16"]
  countries   = ["UZ"]
}
`
	ipNetworkDataSource = `
data "pfptmeta_ip_network" "in" {
  id = pfptmeta_ip_network.in.id
}
`
)

func TestAccResourceIPNetwork(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("content_category", "v1/ip_networks"),
		Steps: []resource.TestStep{
			{
				Config: ipNetworkStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_ip_network.in", "id", regexp.MustCompile("^ipn-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_ip_network.in", "name", "in"),
					resource.TestCheckResourceAttr("pfptmeta_ip_network.in", "description", "in desc"),
					resource.TestCheckResourceAttr("pfptmeta_ip_network.in", "cidrs.0", "0.0.0.0/32"),
					resource.TestCheckResourceAttr("pfptmeta_ip_network.in", "countries.0", "US"),
				),
			},
			{
				Config: ipNetworkStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_ip_network.in", "id", regexp.MustCompile("^ipn-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_ip_network.in", "name", "in 1"),
					resource.TestCheckResourceAttr("pfptmeta_ip_network.in", "description", "in desc 1"),
					resource.TestCheckResourceAttr("pfptmeta_ip_network.in", "cidrs.0", "192.5.0.0/16"),
					resource.TestCheckResourceAttr("pfptmeta_ip_network.in", "countries.0", "UZ"),
				),
			},
		},
	})
}

func TestAccDataSourceIPNetwork(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: ipNetworkStep1 + ipNetworkDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_ip_network.in", "id", regexp.MustCompile("^ipn-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_ip_network.in", "name", "in"),
					resource.TestCheckResourceAttr("pfptmeta_ip_network.in", "description", "in desc"),
					resource.TestCheckResourceAttr("pfptmeta_ip_network.in", "cidrs.0", "0.0.0.0/32"),
					resource.TestCheckResourceAttr("pfptmeta_ip_network.in", "countries.0", "US"),
				),
			},
		},
	})
}
