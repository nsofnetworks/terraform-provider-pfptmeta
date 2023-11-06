package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	scanRuleDependencies = `
resource "pfptmeta_threat_category" "malicious" {
  name             = "Malicious Threat"
  confidence_level = "LOW"
  risk_level       = "LOW"
  countries        = ["IR", "KP"]
  types = [
    "Bitcoin Related", "Blackhole", "Botnets", "Brute Forcer", "CnC", "Compromised", "Drop", "EXE Source",
    "Fake AV", "Keyloggers and Monitoring", "Malware Sites", "Mobile CnC", "Mobile Spyware CnC", "P2P CnC",
    "Phishing and Other Frauds", "Spyware and Adware", "Tor"
  ]
}

resource "pfptmeta_content_category" "cc" {
  name             = "Strict category"
  confidence_level = "LOW"
  types = [
    "Sex Education", "Nudity", "Abused Drugs", "Marijuana", "Swimsuits and Intimate Apparel", "Violence",
    "Gross", "Adult and Pornography", "Weapons", "Hate and Racism", "Gambling"
  ]
  urls = ["espn.com"]
}

data "pfptmeta_catalog_app" "dropbox" {
  name     = "Dropbox"
  category = "Collaboration"
}

resource "pfptmeta_cloud_app" "dropbox_personal" {
  name        = "Dropbox Personal"
  app         = data.pfptmeta_catalog_app.dropbox.id
  tenant_type = "Personal"
}
`
	scanResourceStep1 = `
resource "pfptmeta_scan_rule" "default_rule" {
  name                     = "rule name"
  description              = "rule desc"
  apply_to_org             = true
  action                   = "BLOCK"
  all_supported_file_types = true
  cloud_apps               = [pfptmeta_cloud_app.dropbox_personal.id]
  user_actions             = ["UPLOAD", "DOWNLOAD"]
  priority                 = 15
  filter_expression        = "test:pass"
}
`
	scanResourceStep2 = `
resource "pfptmeta_scan_rule" "default_rule" {
  name                     = "rule name 1"
  description              = "rule desc 1"
  apply_to_org             = true
  action                   = "LOG"
  file_types               = ["arj", "csv", "doc"]
  user_actions             = ["EDIT"]
  priority                 = 20
  cloud_apps               = [pfptmeta_cloud_app.dropbox_personal.id]
  threat_categories        = [pfptmeta_threat_category.malicious.id]
  content_categories       = [pfptmeta_content_category.cc.id]
  filter_expression        = "crwdzta:high"
}
`
	datasourceScanDependencies = `
resource "pfptmeta_content_category" "cc" {
  name             = "for data-source test"
  confidence_level = "LOW"
  types = ["Sex Education"]
}
        
resource "pfptmeta_scan_rule" "default_rule" {
  name                     = "rule for data-source"
  description              = "rule desc"
  apply_to_org             = true
  action                   = "LOG"
  file_types               = ["arj", "csv", "doc"]
  user_actions             = ["EDIT"]
  priority                 = 30
  content_categories       = [pfptmeta_content_category.cc.id]
  filter_expression        = "crwdzta:high"
}
`
	scanForDataSource = `
data "pfptmeta_scan_rule" "scan" {
  id     = pfptmeta_scan_rule.default_rule.id
}
`
)

func TestAccResourceScanRule(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("scan_rule", "v1/scan_rules"),
		Steps: []resource.TestStep{
			{
				Config: scanRuleDependencies + scanResourceStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_scan_rule.default_rule", "id", regexp.MustCompile("^sr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "name", "rule name"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "description", "rule desc"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "action", "BLOCK"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "all_supported_file_types", "true"),
					resource.TestCheckResourceAttrPair("pfptmeta_scan_rule.default_rule", "cloud_apps.0",
						"pfptmeta_cloud_app.dropbox_personal", "id"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "priority", "15"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "user_actions.0", "UPLOAD"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "user_actions.1", "DOWNLOAD"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "filter_expression", "test:pass"),
				),
			},
			{
				Config: scanRuleDependencies + scanResourceStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_scan_rule.default_rule", "id", regexp.MustCompile("^sr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "name", "rule name 1"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "description", "rule desc 1"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "action", "LOG"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "file_types.0", "arj"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "file_types.1", "csv"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "file_types.2", "doc"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "priority", "20"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "user_actions.0", "EDIT"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "user_actions.#", "1"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "filter_expression", "crwdzta:high"),
					resource.TestCheckResourceAttrPair("pfptmeta_scan_rule.default_rule", "cloud_apps.0",
						"pfptmeta_cloud_app.dropbox_personal", "id"),
					resource.TestCheckResourceAttrPair("pfptmeta_scan_rule.default_rule", "threat_categories.0",
						"pfptmeta_threat_category.malicious", "id"),
					resource.TestCheckResourceAttrPair("pfptmeta_scan_rule.default_rule", "content_categories.0",
						"pfptmeta_content_category.cc", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceScanRule(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("scan_rule", "v1/scan_rules"),
		Steps: []resource.TestStep{
			{
				Config: datasourceScanDependencies + scanForDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_scan_rule.scan", "id", regexp.MustCompile("^sr-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_scan_rule.scan", "name", "rule for data-source"),
					resource.TestCheckResourceAttr("data.pfptmeta_scan_rule.scan", "description", "rule desc"),
					resource.TestCheckResourceAttr("data.pfptmeta_scan_rule.scan", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_scan_rule.scan", "action", "LOG"),
					resource.TestCheckResourceAttr("data.pfptmeta_scan_rule.scan", "file_types.0", "arj"),
					resource.TestCheckResourceAttr("data.pfptmeta_scan_rule.scan", "file_types.1", "csv"),
					resource.TestCheckResourceAttr("data.pfptmeta_scan_rule.scan", "file_types.2", "doc"),
					resource.TestCheckResourceAttr("data.pfptmeta_scan_rule.scan", "priority", "30"),
					resource.TestCheckResourceAttr("data.pfptmeta_scan_rule.scan", "user_actions.0", "EDIT"),
					resource.TestCheckResourceAttr("data.pfptmeta_scan_rule.scan", "user_actions.#", "1"),
					resource.TestCheckResourceAttr("data.pfptmeta_scan_rule.scan", "filter_expression", "crwdzta:high"),
					resource.TestCheckResourceAttr("pfptmeta_scan_rule.default_rule", "threat_categories.#", "0"),
					resource.TestCheckResourceAttrPair("pfptmeta_scan_rule.default_rule", "content_categories.0",
						"pfptmeta_content_category.cc", "id"),
				),
			},
		},
	})
}
