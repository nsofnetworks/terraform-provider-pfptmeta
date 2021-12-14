package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceMetaportFailover(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("metaport_failover", "v1/metaport_failovers"),
		Steps: []resource.TestStep{
			{
				Config: testAccMetaportFailoverStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_metaport_failover.failover", "id", regexp.MustCompile("^mpf-.+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "name", "mf-name",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "description", "mf-description",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "failback.0.trigger", "manual",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "failover.0.delay", "1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "failover.0.threshold", "0",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "failover.0.trigger", "auto",
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_metaport_failover.failover", "notification_channels.0", regexp.MustCompile("^nch-.+$"),
					),
				),
			},
			{
				Config: testAccMetaportFailoverStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_metaport_failover.failover", "cluster_1", regexp.MustCompile("^mpc-.+$")),
					resource.TestMatchResourceAttr(
						"pfptmeta_metaport_failover.failover", "cluster_2", regexp.MustCompile("^mpc-.+$")),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "failback.0.trigger", "manual",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "failover.0.delay", "1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "failover.0.threshold", "0",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "failover.0.trigger", "auto",
					),
					resource.TestCheckNoResourceAttr(
						"pfptmeta_metaport_failover.failover", "notification_channels.0",
					),
				),
			},
			{
				Config: testAccMetaportFailoverStep3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_metaport_failover.failover", "mapped_elements.0", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_metaport_failover.failover", "cluster_1", regexp.MustCompile("^mpc-.+$")),
					resource.TestMatchResourceAttr(
						"pfptmeta_metaport_failover.failover", "cluster_2", regexp.MustCompile("^mpc-.+$")),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "failback.0.trigger", "manual",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "failover.0.delay", "1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "failover.0.threshold", "0",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_metaport_failover.failover", "failover.0.trigger", "auto",
					),
				),
			},
			{
				Config: testAccMetaportFailoverStep4,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckNoResourceAttr("pfptmeta_metaport_failover.failover", "mapped_elements"),
					resource.TestCheckNoResourceAttr("pfptmeta_metaport_failover.failover", "cluster_1.#"),
					resource.TestCheckNoResourceAttr("pfptmeta_metaport_failover.failover", "cluster_2.#"),
				),
			},
		},
	})
}

func TestAccDataSourceMetaportFailover(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("metaport_failover", "v1/metaport_failovers"),
		Steps: []resource.TestStep{
			{
				Config: testAccMetaportFailoverStep1 + testAccMetaportFailoverDatasource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_metaport_failover.failover", "id", regexp.MustCompile("^mpf-.+$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_metaport_failover.failover", "name", "mf-name",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_metaport_failover.failover", "description", "mf-description",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_metaport_failover.failover", "failback.0.trigger", "manual",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_metaport_failover.failover", "failover.0.delay", "1",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_metaport_failover.failover", "failover.0.threshold", "0",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_metaport_failover.failover", "failover.0.trigger", "auto",
					),
				),
			},
		},
	})
}

const testAccMetaportFailoverStep1 = `
resource "pfptmeta_notification_channel" "mail" {
  name        = "mail-channel"
  description = "mail channel description"
  email_config {
    recipients = ["user1@example.com", "user2@example.com"]
  }
}

resource "pfptmeta_metaport_failover" "failover" {
  name               = "mf-name"
  description        = "mf-description"
  notification_channels = [pfptmeta_notification_channel.mail.id]
  failback {
	trigger   = "manual"
  }
  failover {
    delay     = 1
    threshold = 0
    trigger   = "auto"
  }
}
`

const testAccMetaportFailoverStep2 = `
resource "pfptmeta_metaport" "metaport1" {
  name = "metaport1"
}

resource "pfptmeta_metaport" "metaport2" {
  name = "metaport2"
}

resource "pfptmeta_metaport_cluster" "metaport_cluster1" {
  name      = "metaport cluster1"
  metaports = [pfptmeta_metaport.metaport1.id]
}

resource "pfptmeta_metaport_cluster" "metaport_cluster2" {
  name      = "metaport cluster2"
  metaports = [pfptmeta_metaport.metaport2.id]
}

resource "pfptmeta_metaport_failover" "failover" {
  name               = "mf-name"
  description        = "mf-description"
  cluster_1          = pfptmeta_metaport_cluster.metaport_cluster1.id
  cluster_2          = pfptmeta_metaport_cluster.metaport_cluster2.id
  failback {
	trigger   = "manual"
  }
  failover {
    delay     = 1
    threshold = 0
    trigger   = "auto"
  }
}
`

const testAccMetaportFailoverStep3 = `
resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "ms"
  mapped_subnets = ["0.0.0.0/0"]
}

resource "pfptmeta_metaport" "metaport1" {
  name = "metaport1"
}

resource "pfptmeta_metaport" "metaport2" {
  name = "metaport2"
}

resource "pfptmeta_metaport_cluster" "metaport_cluster1" {
  name      = "metaport cluster1"
  metaports = [pfptmeta_metaport.metaport1.id]
}

resource "pfptmeta_metaport_cluster" "metaport_cluster2" {
  name      = "metaport cluster2"
  metaports = [pfptmeta_metaport.metaport2.id]
}

resource "pfptmeta_metaport_failover" "failover" {
  name               = "mf-name"
  description        = "mf-description"
  cluster_1          = pfptmeta_metaport_cluster.metaport_cluster1.id
  cluster_2          = pfptmeta_metaport_cluster.metaport_cluster2.id
  mapped_elements    = [pfptmeta_network_element.mapped-subnet.id]
  failback {
	trigger   = "manual"
  }
  failover {
    delay     = 1
    threshold = 0
    trigger   = "auto"
  }
}
`

const testAccMetaportFailoverStep4 = `
resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "ms"
  mapped_subnets = ["0.0.0.0/0"]
}

resource "pfptmeta_metaport" "metaport1" {
  name = "metaport1"
}

resource "pfptmeta_metaport" "metaport2" {
  name = "metaport2"
}

resource "pfptmeta_metaport_cluster" "metaport_cluster1" {
  name      = "metaport cluster1"
  metaports = [pfptmeta_metaport.metaport1.id]
}

resource "pfptmeta_metaport_cluster" "metaport_cluster2" {
  name      = "metaport cluster2"
  metaports = [pfptmeta_metaport.metaport2.id]
}

resource "pfptmeta_metaport_failover" "failover" {
  name               = "mf-name"
  description        = "mf-description"
  failback {
	trigger   = "manual"
  }
  failover {
    delay     = 1
    threshold = 0
    trigger   = "auto"
  }
}
`

const testAccMetaportFailoverDatasource = `

data "pfptmeta_metaport_failover" "failover" {
  id = pfptmeta_metaport_failover.failover.id
}`
