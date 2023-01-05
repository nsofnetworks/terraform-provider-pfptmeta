package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	dlpRuleDependencies = `
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
	dlpResourceStep1 = `
resource "pfptmeta_dlp_rule" "default_rule" {
  name                     = "rule name"
  description              = "rule desc"
  apply_to_org             = true
  action                   = "BLOCK"
  alert_level              = "MEDIUM"
  all_supported_file_types = true
  file_parts               = ["CONTENT"]
  cloud_apps               = [pfptmeta_cloud_app.dropbox_personal.id]
  user_actions             = ["UPLOAD", "DOWNLOAD"]
  priority                 = 15
  filter_expression        = "test:pass"
}
`
	dlpResourceStep2 = `
resource "pfptmeta_dlp_rule" "default_rule" {
  name                     = "rule name 1"
  description              = "rule desc 1"
  apply_to_org             = true
  action                   = "LOG"
  alert_level              = "HIGH"
  file_types               = ["arj", "csv", "doc"]
  file_parts               = ["FILE_NAME"]
  user_actions             = ["EDIT"]
  priority                 = 20
  threat_types             = ["Bitcoin Related"]
}
`
	datasourceDLPDependencies = `
resource "pfptmeta_dlp_rule" "default_rule" {
  name                     = "rule for data-source"
  description              = "rule desc"
  apply_to_org             = true
  action                   = "LOG"
  alert_level              = "HIGH"
  file_types               = ["arj", "csv", "doc"]
  file_parts               = ["FILE_NAME"]
  user_actions             = ["EDIT"]
  priority                 = 30
  threat_types             = ["Bitcoin Related"]
}
`
	dlpForDataSource = `
data "pfptmeta_dlp_rule" "dlp" {
  id     = pfptmeta_dlp_rule.default_rule.id
}
`
)

func TestAccResourceDLPRule(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("dlp_rule", "v1/dlp_rules"),
		Steps: []resource.TestStep{
			{
				Config: dlpRuleDependencies + dlpResourceStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_dlp_rule.default_rule", "id", regexp.MustCompile("^dlp-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "name", "rule name"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "description", "rule desc"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "action", "BLOCK"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "alert_level", "MEDIUM"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "all_supported_file_types", "true"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "file_parts.0", "CONTENT"),
					resource.TestCheckResourceAttrPair("pfptmeta_dlp_rule.default_rule", "cloud_apps.0",
						"pfptmeta_cloud_app.dropbox_personal", "id"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "priority", "15"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "user_actions.0", "UPLOAD"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "user_actions.1", "DOWNLOAD"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "filter_expression", "test:pass"),
				),
			},
			{
				Config: dlpResourceStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_dlp_rule.default_rule", "id", regexp.MustCompile("^dlp-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "name", "rule name 1"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "description", "rule desc 1"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "action", "LOG"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "alert_level", "HIGH"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "file_types.0", "arj"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "file_types.1", "csv"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "file_types.2", "doc"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "file_parts.0", "FILE_NAME"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "priority", "20"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "user_actions.0", "EDIT"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "user_actions.#", "1"),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "filter_expression", ""),
					resource.TestCheckResourceAttr("pfptmeta_dlp_rule.default_rule", "threat_types.0", "Bitcoin Related"),
				),
			},
		},
	})
}

func TestAccDataSourceDLPRule(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("dlp_rule", "v1/dlp_rules"),
		Steps: []resource.TestStep{
			{
				Config: datasourceDLPDependencies + dlpForDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_dlp_rule.dlp", "id", regexp.MustCompile("^dlp-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "name", "rule for data-source"),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "description", "rule desc"),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "action", "LOG"),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "alert_level", "HIGH"),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "file_types.0", "arj"),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "file_types.1", "csv"),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "file_types.2", "doc"),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "file_parts.0", "FILE_NAME"),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "priority", "30"),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "user_actions.0", "EDIT"),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "user_actions.#", "1"),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "filter_expression", ""),
					resource.TestCheckResourceAttr("data.pfptmeta_dlp_rule.dlp", "threat_types.0", "Bitcoin Related"),
				),
			},
		},
	})
}
