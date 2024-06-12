package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	urlFilteringRuleDependencies = `
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

resource "pfptmeta_time_frame" "work_hours" {
  name = "Work Hours"
  days = ["monday", "tuesday", "wednesday", "thursday", "friday"]
  start_time {
    hour   = 8
    minute = 0
  }
  end_time {
    hour   = 18
    minute = 0
  }
}

data "pfptmeta_catalog_app" "salesforce" {
  name     = "Salesforce"
  category = "Business and Finance"
}

resource "pfptmeta_cloud_app" "salesforce" {
  name = "salesforce"
  app  = data.pfptmeta_catalog_app.salesforce.id
  urls = ["my.salesforce.com"]
}

resource "pfptmeta_user" "user" {
  given_name  = "ufr"
  family_name = "user"
  email       = "ufr.user@example.com"
}
`
	ufrResourceStep1 = `
resource "pfptmeta_url_filtering_rule" "default_rule" {
  name                         = "ufr"
  description                  = "ufr desc"
  apply_to_org                 = true
  action                       = "BLOCK"
  advanced_threat_protection   = true
  threat_categories            = [pfptmeta_threat_category.malicious.id]
  forbidden_content_categories = [pfptmeta_content_category.cc.id]
  priority                     = 94
  warn_ttl                     = 15
  filter_expression            = "crwd_agent:fail"
  schedule                     = [pfptmeta_time_frame.work_hours.id]
}

resource "pfptmeta_url_filtering_rule" "high_risk" {
  name              = "ufr 2"
  apply_to_org      = true
  action            = "BLOCK"
  cloud_apps        = [pfptmeta_cloud_app.salesforce.id]
  priority          = 90
  warn_ttl          = 15
}
`
	ufrResourceStep2 = `
resource "pfptmeta_url_filtering_rule" "default_rule" {
  name                         = "ufr 1"
  description                  = "ufr desc 1"
  sources                      = [pfptmeta_user.user.id]
  action                       = "ISOLATE"
  advanced_threat_protection   = false
  forbidden_content_categories = [pfptmeta_content_category.cc.id]
  priority                     = 50
  warn_ttl                     = 15
  filter_expression            = "crwdzta:high"
  access_ids	               = ["-LWPEQeNERIJ-479TKyq3a9sffEWXf6n6yf99M8VWFT"]
}

resource "pfptmeta_url_filtering_rule" "high_risk" {
  name              = "ufr 2 2"
  sources           = [pfptmeta_user.user.id]
  action            = "ISOLATE"
  cloud_apps        = [pfptmeta_cloud_app.salesforce.id]
  priority          = 51
  warn_ttl          = 15
}
`
	datasourceUfrDependencies = `
resource "pfptmeta_content_category" "cc" {
  name             = "for data-source test"
  confidence_level = "LOW"
  types = ["Sex Education"]
}

resource "pfptmeta_url_filtering_rule" "default_rule" {
  name                         = "data source ufr"
  description                  = "data source ufr desc"
  apply_to_org                 = true
  action                       = "ISOLATE"
  advanced_threat_protection   = false
  forbidden_content_categories = [pfptmeta_content_category.cc.id]
  priority                     = 50
  warn_ttl                     = 15
  filter_expression            = "crwdzta:high"
}
`
	ufrForDataSource = `
data "pfptmeta_url_filtering_rule" "ufr" {
  id     = pfptmeta_url_filtering_rule.default_rule.id
}
`
)

func TestAccResourceURLFilteringRule(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("url_filtering_rule", "v1/url_filtering_rules"),
		Steps: []resource.TestStep{
			{
				Config: urlFilteringRuleDependencies + ufrResourceStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_url_filtering_rule.default_rule", "id", regexp.MustCompile("^ufr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "name", "ufr"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "description", "ufr desc"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "action", "BLOCK"),
					resource.TestCheckResourceAttrPair("pfptmeta_url_filtering_rule.default_rule", "threat_categories.0",
						"pfptmeta_threat_category.malicious", "id"),
					resource.TestCheckResourceAttrPair("pfptmeta_url_filtering_rule.default_rule", "forbidden_content_categories.0",
						"pfptmeta_content_category.cc", "id"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "priority", "94"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "warn_ttl", "15"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "filter_expression", "crwd_agent:fail"),
					resource.TestCheckResourceAttrPair("pfptmeta_url_filtering_rule.default_rule", "schedule.0",
						"pfptmeta_time_frame.work_hours", "id"),

					resource.TestMatchResourceAttr("pfptmeta_url_filtering_rule.high_risk", "id", regexp.MustCompile("^ufr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "name", "ufr 2"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "action", "BLOCK"),
					resource.TestCheckResourceAttrPair("pfptmeta_url_filtering_rule.high_risk", "cloud_apps.0",
						"pfptmeta_cloud_app.salesforce", "id"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "priority", "90"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "warn_ttl", "15"),
				),
			},
			{
				Config: urlFilteringRuleDependencies + ufrResourceStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_url_filtering_rule.default_rule", "id", regexp.MustCompile("^ufr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "name", "ufr 1"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "description", "ufr desc 1"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "apply_to_org", "false"),
					resource.TestCheckResourceAttrPair("pfptmeta_url_filtering_rule.default_rule", "sources.0",
						"pfptmeta_user.user", "id"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "action", "ISOLATE"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "threat_categories.#", "0"),
					resource.TestCheckResourceAttrPair("pfptmeta_url_filtering_rule.default_rule", "forbidden_content_categories.0",
						"pfptmeta_content_category.cc", "id"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "priority", "50"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "warn_ttl", "15"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "access_ids.0",
						"-LWPEQeNERIJ-479TKyq3a9sffEWXf6n6yf99M8VWF"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "filter_expression", "crwdzta:high"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "schedule.#", "0"),

					resource.TestMatchResourceAttr("pfptmeta_url_filtering_rule.high_risk", "id", regexp.MustCompile("^ufr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "name", "ufr 2 2"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "apply_to_org", "false"),
					resource.TestCheckResourceAttrPair("pfptmeta_url_filtering_rule.high_risk", "sources.0",
						"pfptmeta_user.user", "id"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "action", "ISOLATE"),
					resource.TestCheckResourceAttrPair("pfptmeta_url_filtering_rule.high_risk", "cloud_apps.0",
						"pfptmeta_cloud_app.salesforce", "id"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "priority", "51"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "warn_ttl", "15"),
				),
			},
		},
	})
}

func TestAccDataSourceURLFilteringRule(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("url_filtering_rule", "v1/url_filtering_rules"),
		Steps: []resource.TestStep{
			{
				Config: urlFilteringRuleDependencies + ufrResourceStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_url_filtering_rule.default_rule", "id", regexp.MustCompile("^ufr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "name", "ufr"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "description", "ufr desc"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "action", "BLOCK"),
					resource.TestCheckResourceAttrPair("pfptmeta_url_filtering_rule.default_rule", "threat_categories.0",
						"pfptmeta_threat_category.malicious", "id"),
					resource.TestCheckResourceAttrPair("pfptmeta_url_filtering_rule.default_rule", "forbidden_content_categories.0",
						"pfptmeta_content_category.cc", "id"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "priority", "94"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "warn_ttl", "15"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.default_rule", "filter_expression", "crwd_agent:fail"),
					resource.TestCheckResourceAttrPair("pfptmeta_url_filtering_rule.default_rule", "schedule.0",
						"pfptmeta_time_frame.work_hours", "id"),

					resource.TestMatchResourceAttr("pfptmeta_url_filtering_rule.high_risk", "id", regexp.MustCompile("^ufr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "name", "ufr 2"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "action", "BLOCK"),
					resource.TestCheckResourceAttrPair("pfptmeta_url_filtering_rule.high_risk", "cloud_apps.0",
						"pfptmeta_cloud_app.salesforce", "id"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "priority", "90"),
					resource.TestCheckResourceAttr("pfptmeta_url_filtering_rule.high_risk", "warn_ttl", "15"),
				),
			},
			{
				Config: datasourceUfrDependencies + ufrForDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_url_filtering_rule.ufr", "id", regexp.MustCompile("^ufr-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_url_filtering_rule.ufr", "name", "data source ufr"),
					resource.TestCheckResourceAttr("data.pfptmeta_url_filtering_rule.ufr", "description", "data source ufr desc"),
					resource.TestCheckResourceAttr("data.pfptmeta_url_filtering_rule.ufr", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_url_filtering_rule.ufr", "action", "ISOLATE"),
					resource.TestCheckResourceAttr("data.pfptmeta_url_filtering_rule.ufr", "advanced_threat_protection", "false"),
					resource.TestCheckResourceAttr("data.pfptmeta_url_filtering_rule.ufr", "threat_categories.#", "0"),
					resource.TestCheckResourceAttrPair("data.pfptmeta_url_filtering_rule.ufr", "forbidden_content_categories.0",
						"pfptmeta_content_category.cc", "id"),
					resource.TestCheckResourceAttr("data.pfptmeta_url_filtering_rule.ufr", "priority", "50"),
					resource.TestCheckResourceAttr("data.pfptmeta_url_filtering_rule.ufr", "warn_ttl", "15"),
					resource.TestCheckResourceAttr("data.pfptmeta_url_filtering_rule.ufr", "filter_expression", "crwdzta:high"),
				),
			},
		},
	})
}
