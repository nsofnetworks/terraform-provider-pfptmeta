package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	egressRouteDependencies = `
resource "pfptmeta_group" "group" {
  name = "egress-rule-group"
}

resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "mapped subnet name"
  mapped_subnets = ["10.20.30.0/24"]
}

data "pfptmeta_location" "new_york" {
  name = "LGA"
}`
	egressRouteWithMappedSubnet = `
resource "pfptmeta_egress_route" "egress" {
  name         = "er-name"
  description  = "er-description"
  sources      = [pfptmeta_group.group.id]
  destinations = ["example.com"]
  via          = pfptmeta_network_element.mapped-subnet.id
}`
	egressRouteWithLocation = `
resource "pfptmeta_egress_route" "egress" {
  name         = "er-name1"
  sources      = [pfptmeta_group.group.id]
  destinations = ["example1.com", "example2.com"]
  via          = data.pfptmeta_location.new_york.name
}`
	egressRouteWithDirect = `
resource "pfptmeta_egress_route" "egress" {
  name         = "er-name"
  sources      = [pfptmeta_group.group.id]
  destinations = ["example3.com"]
  via          = "DIRECT"
}`
	dataSourceEgressRoute = `
data "pfptmeta_egress_route" "egress" {
id = pfptmeta_egress_route.egress.id
}`
)

func TestAccResourceEgressRoute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("group", "v1/groups"),
		Steps: []resource.TestStep{
			{
				Config: egressRouteDependencies + egressRouteWithMappedSubnet,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_egress_route.egress", "id", regexp.MustCompile("^er-.*$")),
					resource.TestCheckResourceAttr("pfptmeta_egress_route.egress", "name", "er-name"),
					resource.TestCheckResourceAttr("pfptmeta_egress_route.egress", "description", "er-description"),
					resource.TestMatchResourceAttr("pfptmeta_egress_route.egress", "sources.0", regexp.MustCompile("^grp-.*$")),
					resource.TestCheckResourceAttr("pfptmeta_egress_route.egress", "destinations.0", "example.com"),
					resource.TestMatchResourceAttr("pfptmeta_egress_route.egress", "via", regexp.MustCompile("^ne-.*$")),
				),
			},
			{
				Config: egressRouteDependencies + egressRouteWithLocation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_egress_route.egress", "id", regexp.MustCompile("^er-.*$")),
					resource.TestCheckResourceAttr("pfptmeta_egress_route.egress", "name", "er-name1"),
					resource.TestCheckResourceAttr("pfptmeta_egress_route.egress", "description", ""),
					resource.TestMatchResourceAttr("pfptmeta_egress_route.egress", "sources.0", regexp.MustCompile("^grp-.*$")),
					resource.TestCheckResourceAttr("pfptmeta_egress_route.egress", "destinations.0", "example1.com"),
					resource.TestCheckResourceAttr("pfptmeta_egress_route.egress", "destinations.1", "example2.com"),
					resource.TestCheckResourceAttr("pfptmeta_egress_route.egress", "via", "LGA"),
				),
			},
			{
				Config: egressRouteDependencies + egressRouteWithDirect,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_egress_route.egress", "id", regexp.MustCompile("^er-.*$")),
					resource.TestCheckResourceAttr("pfptmeta_egress_route.egress", "name", "er-name"),
					resource.TestCheckResourceAttr("pfptmeta_egress_route.egress", "description", ""),
					resource.TestMatchResourceAttr("pfptmeta_egress_route.egress", "sources.0", regexp.MustCompile("^grp-.*$")),
					resource.TestCheckResourceAttr("pfptmeta_egress_route.egress", "destinations.0", "example3.com"),
					resource.TestCheckResourceAttr("pfptmeta_egress_route.egress", "via", "DIRECT"),
				),
			},
		},
	})
}

func TestAccDataSourceEgressRoute(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: egressRouteDependencies + egressRouteWithMappedSubnet + dataSourceEgressRoute,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_egress_route.egress", "id", regexp.MustCompile("^er-.*$")),
					resource.TestCheckResourceAttr("data.pfptmeta_egress_route.egress", "name", "er-name"),
					resource.TestCheckResourceAttr("data.pfptmeta_egress_route.egress", "description", "er-description"),
					resource.TestMatchResourceAttr("data.pfptmeta_egress_route.egress", "sources.0", regexp.MustCompile("^grp-.*$")),
					resource.TestCheckResourceAttr("data.pfptmeta_egress_route.egress", "destinations.0", "example.com"),
					resource.TestMatchResourceAttr("data.pfptmeta_egress_route.egress", "via", regexp.MustCompile("^ne-.*$")),
				),
			},
		},
	})
}
