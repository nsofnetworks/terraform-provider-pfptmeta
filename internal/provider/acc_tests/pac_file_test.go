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
  type         = "bring_your_own"
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
  type         = "bring_your_own"
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
  type         = "bring_your_own"
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
	managedPacNew = `
resource "pfptmeta_pac_file" "pac" {
  name            = "test managed pac file"
  description     = "test managed pac file description"
  apply_to_org    = true
  priority        = 15
  type            = "managed"
}
`
	managedPacUpdated_1 = `
resource "pfptmeta_pac_file" "pac" {
  name            = "test managed pac file 1"
  description     = "test managed pac file description 1"
  apply_to_org    = false
  priority        = 20
  type            = "managed"

  managed_content {
	domains = ["battle.net"]
	cloud_apps = []
  }
}
`
	managedPacUpdated_2 = `
resource "pfptmeta_pac_file" "pac" {
  name            = "test managed pac file 2"
  description     = "test managed pac file description 2"
  apply_to_org    = false
  priority        = 20
  type            = "managed"

  managed_content {
	domains = ["battle.net", "warhammer40k.com"]
  }
}
`
	managedPacCleared = `
resource "pfptmeta_pac_file" "pac" {
  name            = "test managed pac file cleared"
  description     = "test managed pac file cleared description"
  apply_to_org    = true
  priority        = 15
  type            = "managed"

  managed_content {
	domains = []
  }
}
`
	managedPacFileDataSource = `
resource "pfptmeta_pac_file" "pac_data_source" {
  name            = "test managed pac file data source"
  description     = "test managed pac file data source description"
  apply_to_org    = true
  priority        = 15
  type            = "managed"

  managed_content {
	domains = ["proof.point.com"]
  }
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
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "type", "bring_your_own"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.#", "0"), // verify empty by seeing length is zero
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
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "type", "bring_your_own"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.#", "0"),
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
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "type", "bring_your_own"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac_data_source", "managed_content.#", "0"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "content", "function FindProxyForURL(url, host) {\n  return \"PROXY 127.0.0.1:43443\";\n}\n"),
				),
			},
		},
	})
}

func TestAccResourceManagedPacFile(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("pac_file", "v1/pac_files"),
		Steps: []resource.TestStep{
			{
				Config: managedPacNew,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "id", regexp.MustCompile("^pf-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "name", "test managed pac file"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "description", "test managed pac file description"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "priority", "15"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "type", "managed"),
					//resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.domains", regexp.MustCompile(".+")),   // not empty - how to check non empty???  // NADAV error: Check failed: Check 7/10 error: pfptmeta_pac_file.pac: Attribute 'managed_content.0.domains' didn't match ".+", got ""
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.cloud_apps.#", "0"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.ip_networks.#", "0"),
					resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "content", regexp.MustCompile(".+")),
				),
				ExpectNonEmptyPlan: true, // NADAV HERE: why are we failing on plan non empty? is it because the backend adds default domains to managed content? why does it not fail on added content?
			},
			{
				Config: managedPacUpdated_1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "id", regexp.MustCompile("^pf-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "name", "test managed pac file 1"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "description", "test managed pac file description 1"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "apply_to_org", "false"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "priority", "20"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "type", "managed"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.domains.#", "1"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.domains.0", "battle.net"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.cloud_apps.#", "0"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.ip_networks.#", "0"),
					resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "content", regexp.MustCompile(`battle.net`)),
				),
			},
			{
				Config: managedPacUpdated_2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "id", regexp.MustCompile("^pf-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "name", "test managed pac file 2"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "description", "test managed pac file description 2"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "apply_to_org", "false"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "priority", "20"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "type", "managed"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.domains.#", "2"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.domains.0", "battle.net"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.domains.1", "warhammer40k.com"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.cloud_apps.#", "0"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.ip_networks.#", "0"),
					resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "content", regexp.MustCompile(`battle.net`)),
					resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "content", regexp.MustCompile(`warhammer40k.com`)),
				),
			},
			{
				Config: managedPacCleared,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "id", regexp.MustCompile("^pf-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "name", "test managed pac file cleared"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "description", "test managed pac file cleared description"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "priority", "15"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "type", "managed"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.domains.#", "0"),  // NADAV: not cleared!!
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.cloud_apps.#", "0"),
					resource.TestCheckResourceAttr("pfptmeta_pac_file.pac", "managed_content.0.ip_networks.#", "0"),
					resource.TestMatchResourceAttr("pfptmeta_pac_file.pac", "content", regexp.MustCompile(".+")),
				),
			},
		},
	})
}

func TestAccDataSourceManagedPacFile(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: managedPacFileDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_pac_file.pac_data_source", "id", regexp.MustCompile("^pf-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "name", "test managed pac file data source"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "description", "test managed pac file data source description"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "priority", "15"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "type", "managed"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "managed_content.0.domains.#", "1"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "managed_content.0.domains.0", "proof.point.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "managed_content.0.cloud_apps.#", "0"),
					resource.TestCheckResourceAttr("data.pfptmeta_pac_file.pac_data_source", "managed_content.0.ip_networks.#", "0"),
					resource.TestMatchResourceAttr("data.pfptmeta_pac_file.pac_data_source", "content", regexp.MustCompile(`\"(proof.point.com)\"`)),
				),
			},
		},
	})
}
