package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	threatCategoryStep1 = `
resource "pfptmeta_threat_category" "tc" {
  name             = "tc"
  description      = "tc desc"
  confidence_level = "MEDIUM"
  risk_level       = "HIGH"
  countries        = ["EG"]
  types            = ["Scanner"]
}
`
	threatCategoryStep2 = `
resource "pfptmeta_threat_category" "tc" {
  name             = "tc1"
  description      = "tc desc 1"
  confidence_level = "HIGH"
  risk_level       = "MEDIUM"
  countries        = ["AF", "EG"]
  types            = ["Peer to Peer", "Scanner"]
}
`
	threatCategoryDataSource = `
data "pfptmeta_threat_category" "tc" {
  id = pfptmeta_threat_category.tc.id
}
`
)

func TestAccResourceThreatCategory(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("threat_category", "v1/threat_categories"),
		Steps: []resource.TestStep{
			{
				Config: threatCategoryStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_threat_category.tc", "id", regexp.MustCompile("^tc-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "name", "tc"),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "description", "tc desc"),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "confidence_level", "MEDIUM"),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "risk_level", "HIGH"),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "countries.0", "EG"),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "types.0", "Scanner"),
				),
			},
			{
				Config: threatCategoryStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_threat_category.tc", "id", regexp.MustCompile("^tc-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "name", "tc1"),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "description", "tc desc 1"),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "confidence_level", "HIGH"),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "risk_level", "MEDIUM"),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "countries.0", "AF"),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "countries.1", "EG"),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "types.0", "Peer to Peer"),
					resource.TestCheckResourceAttr("pfptmeta_threat_category.tc", "types.1", "Scanner"),
				),
			},
		},
	})
}

func TestAccDataSourceThreatCategory(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: threatCategoryStep1 + threatCategoryDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_threat_category.tc", "id", regexp.MustCompile("^tc-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_threat_category.tc", "name", "tc"),
					resource.TestCheckResourceAttr("data.pfptmeta_threat_category.tc", "description", "tc desc"),
					resource.TestCheckResourceAttr("data.pfptmeta_threat_category.tc", "confidence_level", "MEDIUM"),
					resource.TestCheckResourceAttr("data.pfptmeta_threat_category.tc", "risk_level", "HIGH"),
					resource.TestCheckResourceAttr("data.pfptmeta_threat_category.tc", "countries.0", "EG"),
					resource.TestCheckResourceAttr("data.pfptmeta_threat_category.tc", "types.0", "Scanner"),
				),
			},
		},
	})
}
