package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	userSettingsGroup = `
resource "pfptmeta_group" "ds_group" {
  name = "user settings group"
}
`
	userSettingsResourceStep1 = `
resource "pfptmeta_user_settings" "settings" {
  name                 = "settings-name"
  description          = "settings-desc"
  apply_on_org         = true
  proxy_pops           = "POPS_WITH_DEDICATED_IPS"
  max_devices_per_user = "5"
  prohibited_os        = ["macOS", "iOS"]
  sso_mandatory        = true
  mfa_required         = true
  allowed_factors      = ["SMS"]
  password_expiration  = 30
}
`
	userSettingsResourceStep2 = `
resource "pfptmeta_user_settings" "settings" {
  name                 = "settings-name1"
  description          = "settings-desc1"
  apply_to_entities    = [pfptmeta_group.ds_group.id]
  proxy_pops           = "ALL_POPS"
  prohibited_os        = ["Windows", "ChromeOS"]
  sso_mandatory        = false
  password_expiration  = 15
}
`
	userSettingsDataSource = `
data "pfptmeta_user_settings" "settings" {
  id = pfptmeta_user_settings.settings.id
}`
)

func TestAccResourceUserSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("user_settings", "v1/settings/user"),
		Steps: []resource.TestStep{
			{
				Config: userSettingsResourceStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_user_settings.settings", "id", regexp.MustCompile("^as-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "name", "settings-name"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "description", "settings-desc"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "apply_on_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "proxy_pops", "POPS_WITH_DEDICATED_IPS"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "max_devices_per_user", "5"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "prohibited_os.0", "macOS"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "prohibited_os.1", "iOS"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "sso_mandatory", "true"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "mfa_required", "true"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "allowed_factors.0", "SMS"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "password_expiration", "30"),
				),
			},
			{
				Config: userSettingsGroup + userSettingsResourceStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_user_settings.settings", "id", regexp.MustCompile("^as-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "name", "settings-name1"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "description", "settings-desc1"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "apply_on_org", "false"),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_user_settings.settings", "apply_to_entities.0",
						"pfptmeta_group.ds_group", "id"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "proxy_pops", "ALL_POPS"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "max_devices_per_user", ""),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "prohibited_os.0", "Windows"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "prohibited_os.1", "ChromeOS"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "sso_mandatory", "false"),
					resource.TestCheckNoResourceAttr("pfptmeta_user_settings.settings", "mfa_required"),
					resource.TestCheckNoResourceAttr("pfptmeta_user_settings.settings", "allowed_factors"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "password_expiration", "15"),
				),
			},
		},
	})
}

func TestAccDataSourceUserSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("user_settings", "v1/settings/user"),
		Steps: []resource.TestStep{
			{
				Config: userSettingsResourceStep1 + userSettingsDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_user_settings.settings", "id", regexp.MustCompile("^as-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "name", "settings-name"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "description", "settings-desc"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "apply_on_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "proxy_pops", "POPS_WITH_DEDICATED_IPS"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "max_devices_per_user", "5"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "prohibited_os.0", "macOS"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "prohibited_os.1", "iOS"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "sso_mandatory", "true"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "mfa_required", "true"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "allowed_factors.0", "SMS"),
					resource.TestCheckResourceAttr("pfptmeta_user_settings.settings", "password_expiration", "30"),
				),
			},
		},
	})
}
