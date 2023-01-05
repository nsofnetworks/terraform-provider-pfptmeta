package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	fsrDependencies = `
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
	fsrResourceStep1 = `
resource "pfptmeta_file_scanning_rule" "file_scanning" {
  name              = "File scanning"
  description       = "File scanning rule"
  apply_to_org      = true
  block_file_types  = ["exe"]
  malware           = "DOWNLOAD"
  priority          = 15
  filter_expression = "test:pass"
}
`
	fsrResourceStep2 = `
resource "pfptmeta_file_scanning_rule" "file_scanning" {
  name                    = "File scanning 1"
  description             = "File scanning rule 1"
  apply_to_org            = true
  block_file_types        = ["msi"]
  block_content_types     = ["Abortion"]
  block_countries         = ["AD", "AE"]
  block_threat_types      = ["VPN"]
  block_unsupported_files = true
  malware                 = "ALL"
  priority                = 16
  cloud_apps              = [pfptmeta_cloud_app.dropbox_personal.id]
}
`
	datasourceFSRDependencies = `
resource "pfptmeta_file_scanning_rule" "file_scanning" {
  name             = "data-source rule"
  description      = "File scanning rule"
  apply_to_org     = true
  block_file_types = ["exe"]
  malware          = "DOWNLOAD"
  priority         = 17
}
`
	fsrForDataSource = `
data "pfptmeta_file_scanning_rule" "fsr" {
  id     = pfptmeta_file_scanning_rule.file_scanning.id
}
`
)

func TestAccResourceFSR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("file_scanning_rule", "v1/file_scanning_rules"),
		Steps: []resource.TestStep{
			{
				Config: fsrResourceStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "id", regexp.MustCompile("^fsr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "name", "File scanning"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "description", "File scanning rule"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "malware", "DOWNLOAD"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "block_file_types.0", "exe"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "priority", "15"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "filter_expression", "test:pass"),
				),
			},
			{
				Config: fsrDependencies + fsrResourceStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "id", regexp.MustCompile("^fsr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "name", "File scanning 1"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "description", "File scanning rule 1"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "block_file_types.0", "msi"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "block_content_types.0", "Abortion"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "block_countries.0", "AD"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "block_countries.1", "AE"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "block_threat_types.0", "VPN"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "block_unsupported_files", "true"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "malware", "ALL"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "priority", "16"),
					resource.TestCheckResourceAttrPair("pfptmeta_file_scanning_rule.file_scanning", "cloud_apps.0",
						"pfptmeta_cloud_app.dropbox_personal", "id"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "filter_expression", ""),
				),
			},
		},
	})
}

func TestAccDataSourceFSR(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("file_scanning_rule", "v1/file_scanning_rules"),
		Steps: []resource.TestStep{
			{
				Config: datasourceFSRDependencies + fsrForDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "id", regexp.MustCompile("^fsr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "name", "data-source rule"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "description", "File scanning rule"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "malware", "DOWNLOAD"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "block_file_types.0", "exe"),
					resource.TestCheckResourceAttr("pfptmeta_file_scanning_rule.file_scanning", "priority", "17"),
				),
			},
		},
	})
}
