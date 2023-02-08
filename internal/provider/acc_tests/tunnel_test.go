package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceTunnel(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("tunnel", "v1/tunnels"),
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_tunnel.tunnel", "id", regexp.MustCompile("^tun-.+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_tunnel.tunnel", "name", "tunnel name1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_tunnel.tunnel", "description", "tunnel description1",
					),
				),
			},
			{
				Config: testAccTunnelStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pfptmeta_tunnel.tunnel", "name", "tunnel name2",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_tunnel.tunnel", "description", "tunnel description2",
					),
				),
			},
			{
				Config: testAccTunnelStep3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pfptmeta_tunnel.tunnel", "name", "tunnel name3",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_tunnel.tunnel", "description", "tunnel description3",
					),
				),
			},
			{
				Config: testAccTunnelStep4,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pfptmeta_tunnel.tunnel", "name", "tunnel name4",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_tunnel.tunnel", "description", "tunnel description4",
					),
				),
			},
			{
				Config: testAccTunnelStep5,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pfptmeta_tunnel.tunnel", "name", "tunnel name5",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_tunnel.tunnel", "description", "tunnel description5",
					),
				),
			},
		},
	})
}

func TestAccDataSourceTunnel(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("tunnel", "v1/tunnels"),
		Steps: []resource.TestStep{
			{
				Config: testAccTunnelStep1 + testAccTunnelByIDDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_tunnel.tunnel2", "id", regexp.MustCompile("^tun-.+$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_tunnel.tunnel2", "name", "tunnel name1",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_tunnel.tunnel2", "description", "tunnel description1",
					),
				),
			},
			{
				Config: testAccTunnelStep1 + testAccTunnelByNameDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_tunnel.tunnel3", "id", regexp.MustCompile("^tun-.+$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_tunnel.tunnel3", "name", "tunnel name1",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_tunnel.tunnel3", "description", "tunnel description1",
					),
				),
			},
		},
	})
}

const testAccTunnelStep1 = `
resource "pfptmeta_tunnel" "tunnel" {
  name                  = "tunnel name1"
  description           = "tunnel description1"
  gre_config {
      source_ips = ["1.2.3.4", "2.3.4.5"]
  }
}
`

const testAccTunnelStep2 = `
resource "pfptmeta_tunnel" "tunnel" {
  name                  = "tunnel name2"
  description           = "tunnel description2"
  gre_config {
      source_ips = ["1.2.3.4"]
  }
}
`

const testAccTunnelStep3 = `
resource "pfptmeta_tunnel" "tunnel" {
  name                  = "tunnel name3"
  description           = "tunnel description3"
  gre_config {
      source_ips = ["1.2.3.4", "3.4.5.6"]
  }
}
`

const testAccTunnelStep4 = `
resource "pfptmeta_tunnel" "tunnel" {
  name                  = "tunnel name4"
  description           = "tunnel description4"
  gre_config {
      source_ips = ["4.5.6.7"]
  }
}
`

const testAccTunnelStep5 = `
resource "pfptmeta_tunnel" "tunnel" {
  name                  = "tunnel name5"
  description           = "tunnel description5"
}
`

const testAccTunnelByIDDataSource = `

data "pfptmeta_tunnel" "tunnel2" {
  id = pfptmeta_tunnel.tunnel.id
}`

const testAccTunnelByNameDataSource = `

data "pfptmeta_tunnel" "tunnel3" {
  name = "tunnel name1"
}`
