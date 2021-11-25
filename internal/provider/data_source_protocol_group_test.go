package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceProtocolGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceProtocolGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.pfptmeta_protocol_group.HTTP", "id", "pg-NKMzUnJzalxWZKe",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_protocol_group.HTTP", "name", "HTTP",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_protocol_group.HTTP", "protocols.0.from_port", "80",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_protocol_group.HTTP", "protocols.0.to_port", "80",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_protocol_group.HTTP", "protocols.0.proto", "tcp",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_protocol_group.HTTPS", "id", "pg-D6vrsnKbj7q5Kxl",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_protocol_group.HTTPS", "name", "HTTPS",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_protocol_group.HTTPS", "protocols.0.from_port", "443",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_protocol_group.HTTPS", "protocols.0.to_port", "443",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_protocol_group.HTTPS", "protocols.0.proto", "tcp",
					),
				),
			},
		},
	})
}

const testAccDataSourceProtocolGroup = `
data "pfptmeta_protocol_group" "HTTP" {
  id = "pg-NKMzUnJzalxWZKe"
}

data "pfptmeta_protocol_group" "HTTPS" {
  name = "HTTPS"
}
`
