package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMetaportCluster(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("metaport_cluster", "v1/metaport_clusters"),
		Steps: []resource.TestStep{
			{
				Config: testAccMetaportClusterStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_metaport_cluster.metaport_cluster", "id", regexp.MustCompile("^mpc-.+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_cluster.metaport_cluster", "name", "metaport cluster name",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_cluster.metaport_cluster", "description", "metaport cluster description",
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_metaport_cluster.metaport_cluster", "mapped_elements.0", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_metaport_cluster.metaport_cluster", "metaports.0", regexp.MustCompile("^mp-[\\d]+$"),
					),
				),
			},
			{
				Config: testAccMetaportClusterStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_cluster.metaport_cluster", "name", "metaport cluster name1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_cluster.metaport_cluster", "description", "metaport cluster description1",
					),
					resource.TestCheckNoResourceAttr("pfptmeta_metaport_cluster.metaport_cluster", "mapped_elements"),
					resource.TestCheckNoResourceAttr("pfptmeta_metaport_cluster.metaport_cluster", "metaports"),
				),
			},
		},
	})
}

func TestAccDataMetaportCluster(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("metaport_cluster", "v1/metaport_clusters"),
		Steps: []resource.TestStep{
			{
				Config: testAccMetaportClusterStep1 + testAccDataSourceMetaportClusterByID,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_metaport_cluster.metaport_cluster", "id", regexp.MustCompile("^mpc-.+$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_metaport_cluster.metaport_cluster", "name", "metaport cluster name",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_metaport_cluster.metaport_cluster", "description", "metaport cluster description",
					),
					resource.TestMatchResourceAttr(
						"data.pfptmeta_metaport_cluster.metaport_cluster", "mapped_elements.0", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestMatchResourceAttr(
						"data.pfptmeta_metaport_cluster.metaport_cluster", "metaports.0", regexp.MustCompile("^mp-[\\d]+$"),
					),
				),
			},
			{
				Config: testAccMetaportClusterStep1 + testAccDataSourceMetaportClusterByName,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_metaport_cluster.metaport_cluster", "id", regexp.MustCompile("^mpc-.+$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_metaport_cluster.metaport_cluster", "name", "metaport cluster name",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_metaport_cluster.metaport_cluster", "description", "metaport cluster description",
					),
					resource.TestMatchResourceAttr(
						"data.pfptmeta_metaport_cluster.metaport_cluster", "mapped_elements.0", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestMatchResourceAttr(
						"data.pfptmeta_metaport_cluster.metaport_cluster", "metaports.0", regexp.MustCompile("^mp-[\\d]+$"),
					),
				),
			},
		},
	})
}

const testAccMetaportClusterStep1 = `
resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "ms"
  mapped_subnets = ["0.0.0.0/0"]
}

resource "pfptmeta_metaport" "metaport" {
  name                  = "metaport name"
  description           = "metaport description"
}

resource "pfptmeta_metaport_cluster" "metaport_cluster" {
  name = "metaport cluster name"
  description = "metaport cluster description"
  metaports = [pfptmeta_metaport.metaport.id]
  mapped_elements = [pfptmeta_network_element.mapped-subnet.id]
}
`

const testAccMetaportClusterStep2 = `
resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "ms"
  mapped_subnets = ["0.0.0.0/0"]
}

resource "pfptmeta_metaport" "metaport" {
  name                  = "metaport name"
  description           = "metaport description"
}

resource "pfptmeta_metaport_cluster" "metaport_cluster" {
  name = "metaport cluster name1"
  description = "metaport cluster description1"
}
`

const testAccDataSourceMetaportClusterByID = `

data "pfptmeta_metaport_cluster" "metaport_cluster" {
  id = pfptmeta_metaport_cluster.metaport_cluster.id
}`

const testAccDataSourceMetaportClusterByName = `

data "pfptmeta_metaport_cluster" "metaport_cluster" {
  name = "metaport cluster name"
}`
