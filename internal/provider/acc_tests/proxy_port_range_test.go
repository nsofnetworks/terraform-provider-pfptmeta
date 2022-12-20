package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	proxyPortRangeStep1 = `
resource "pfptmeta_proxy_port_range" "ppr" {
  name        = "ppr-plugin-tf1"
  description = "ppr1 desc"
  proto       = "HTTP"
  from_port   = "80"
  to_port     = "85"
}
`
	proxyPortRangeStep2 = `
resource "pfptmeta_proxy_port_range" "ppr" {
  name        = "ppr-tf2"
  description = "ppr2 desc"
  proto       = "HTTPS"
  from_port   = "443"
  to_port     = "445"
}
`

	proxyPortRangeDataSource = `
data "pfptmeta_proxy_port_range" "ppr" {
  id = pfptmeta_proxy_port_range.ppr.id
}
`
)

func TestAccResourceproxyPortRange(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("proxy_port_range", "v1/proxy_port_ranges"),
		Steps: []resource.TestStep{
			{
				Config: proxyPortRangeStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_proxy_port_range.ppr", "id", regexp.MustCompile("^ppr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_proxy_port_range.ppr", "name", "ppr-plugin-tf1"),
					resource.TestCheckResourceAttr("pfptmeta_proxy_port_range.ppr", "description", "ppr1 desc"),
					resource.TestCheckResourceAttr("pfptmeta_proxy_port_range.ppr", "proto", "HTTP"),
					resource.TestCheckResourceAttr("pfptmeta_proxy_port_range.ppr", "from_port", "80"),
					resource.TestCheckResourceAttr("pfptmeta_proxy_port_range.ppr", "to_port", "85"),
				),
			},
			{
				Config: proxyPortRangeStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_proxy_port_range.ppr", "id", regexp.MustCompile("^ppr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_proxy_port_range.ppr", "name", "ppr-tf2"),
					resource.TestCheckResourceAttr("pfptmeta_proxy_port_range.ppr", "description", "ppr2 desc"),
					resource.TestCheckResourceAttr("pfptmeta_proxy_port_range.ppr", "proto", "HTTPS"),
					resource.TestCheckResourceAttr("pfptmeta_proxy_port_range.ppr", "from_port", "443"),
					resource.TestCheckResourceAttr("pfptmeta_proxy_port_range.ppr", "to_port", "445"),
				),
			},
		},
	})
}

func TestAccDataSourceproxyPortRange(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: proxyPortRangeStep1 + proxyPortRangeDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_proxy_port_range.ppr", "id", regexp.MustCompile("^ppr-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_proxy_port_range.ppr", "name", "ppr-plugin-tf1"),
					resource.TestCheckResourceAttr("data.pfptmeta_proxy_port_range.ppr", "description", "ppr1 desc"),
					resource.TestCheckResourceAttr("data.pfptmeta_proxy_port_range.ppr", "proto", "HTTP"),
					resource.TestCheckResourceAttr("data.pfptmeta_proxy_port_range.ppr", "from_port", "80"),
					resource.TestCheckResourceAttr("data.pfptmeta_proxy_port_range.ppr", "to_port", "85"),
				),
			},
		},
	})
}
