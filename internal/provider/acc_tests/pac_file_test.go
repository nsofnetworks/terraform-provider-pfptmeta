package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	byoPacNew = `
resource "pfptmeta_pac_file" "pac" {
  name         = "test byo pac file"
  description  = "test byo pac file description"
  apply_to_org = true
  priority     = 15
  content      = <<EOF
function FindProxyForURL(url, host) {
  return "PROXY 127.0.0.1:43443321";
}
EOF
}
`
	byoPacUpdated = `
resource "pfptmeta_pac_file" "pac" {
  name         = "test byo pac file 1"
  description  = "test byo pac file description 1"
  apply_to_org = false
  priority     = 20
  content      = <<EOF
function FindProxyForURL(url, host) {
  return "PROXY 127.0.0.1:43443";
}
EOF
}
`
	byoPacFileDataSource = `
resource "pfptmeta_pac_file" "pac_data_source" {
  name         = "test byo pac file data source"
  description  = "test byo pac file data source description"
  apply_to_org = true
  priority     = 15
  content      = <<EOF
function FindProxyForURL(url, host) {
  return "PROXY 127.0.0.1:43443";
}
EOF
}

data "pfptmeta_pac_file" "pac_data_source" {
  id = pfptmeta_pac_file.pac_data_source.id
}
`
)

func TestAccResourceBYOPacFile(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("pac_file", "v1/pac_files"),
		Steps: []resource.TestStep{
			{
				Config: byoPacNew,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "id", regexp.MustCompile("^pf-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "name", "test byo pac file"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "description", "test byo pac file description"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "priority", "15"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "content", "function FindProxyForURL(url, host) {\n  return \"PROXY 127.0.0.1:43443321\";\n}\n"),
				),
			},
			{
				Config: byoPacUpdated,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "id", regexp.MustCompile("^pf-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "name", "test byo pac file 1"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "description", "test byo pac file description 1"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "apply_to_org", "false"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "priority", "20"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "content", "function FindProxyForURL(url, host) {\n  return \"PROXY 127.0.0.1:43443\";\n}\n"),
				),
			},
		},
	})
}

func TestAccDataSourceBYOPacFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: byoPacFileDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_pac_file.pac_data_source", "id", regexp.MustCompile("^pf-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "name", "test byo pac file data source"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "description", "test byo pac file data source description"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "priority", "15"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "content", "function FindProxyForURL(url, host) {\n  return \"PROXY 127.0.0.1:43443\";\n}\n"),
				),
			},
		},
	})
}

