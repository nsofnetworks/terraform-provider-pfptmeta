package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	pacFileStep1 = `
resource "pfptmeta_pac_file" "pac" {
  name         = "pac file"
  description  = "pac file description"
  apply_to_org = true
  priority     = 15
  content      = <<EOF
function FindProxyForURL(url, host) {
  return "PROXY 127.0.0.1:43443321";
}
EOF
}
`
	pacFileStep2 = `
resource "pfptmeta_pac_file" "pac" {
  name         = "pac file 1"
  description  = "pac file description 1"
  apply_to_org = false
  priority     = 20
  content      = <<EOF
function FindProxyForURL(url, host) {
  return "PROXY 127.0.0.1:43443";
}
EOF
}
`
	pacFileStepDataSource = `

resource "pfptmeta_pac_file" "pac" {
  name         = "pac file"
  description  = "pac file description"
  apply_to_org = true
  priority     = 15
  content      = <<EOF
function FindProxyForURL(url, host) {
  return "PROXY 127.0.0.1:43443";
}
EOF
}

data "pfptmeta_pac_file" "pac" {
  id = pfptmeta_pac_file.pac.id
}
`
)

func TestAccResourcePacFile(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("pac_file", "v1/pac_files"),
		Steps: []resource.TestStep{
			{
				Config: pacFileStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "id", regexp.MustCompile("^pf-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "name", "pac file"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "description", "pac file description"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "priority", "15"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "content", "function FindProxyForURL(url, host) {\n  return \"PROXY 127.0.0.1:43443321\";\n}\n"),
				),
			},
			{
				Config: pacFileStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "id", regexp.MustCompile("^pf-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "name", "pac file 1"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "description", "pac file description 1"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "apply_to_org", "false"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "priority", "20"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "content", "function FindProxyForURL(url, host) {\n  return \"PROXY 127.0.0.1:43443\";\n}\n"),
				),
			},
		},
	})
}

func TestAccDataSourcePacFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: pacFileStepDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_pac_file.pac", "id", regexp.MustCompile("^pf-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac", "name", "pac file"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac", "description", "pac file description"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac", "priority", "15"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac", "content", "function FindProxyForURL(url, host) {\n  return \"PROXY 127.0.0.1:43443\";\n}\n"),
				),
			},
		},
	})
}
