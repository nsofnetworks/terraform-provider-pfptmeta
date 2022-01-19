package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	deviceSettingsGroup = `
resource "pfptmeta_group" "ds_group" {
  name = "device settings group"
}
`
	deviceSettingsResourceStep1 = `
resource "pfptmeta_device_settings" "settings" {
  name                        = "settings-name"
  description                 = "settings-desc"
  apply_on_org                = true
  tunnel_mode                 = "full"
  ztna_always_on              = true
  proxy_always_on             = true
  overlay_mfa_required        = true
  overlay_mfa_refresh_period  = 20
  vpn_login_browser           = "EXTERNAL"
  session_lifetime            = 15
  session_lifetime_grace      = 3
  protocol_selection_lifetime = 10
  search_domains              = ["example1.com", "example2.com"]
}
`
	deviceSettingsResourceStep2 = `
resource "pfptmeta_device_settings" "settings" {
  name                        = "settings-name1"
  description                 = "settings-desc1"
  apply_to_entities           = [pfptmeta_group.ds_group.id]
  tunnel_mode                 = "split"
  vpn_login_browser           = "USER_DEFINED"
  session_lifetime            = 14
  session_lifetime_grace      = 4
  protocol_selection_lifetime = 12
}
`
	deviceSettingsDataSource = `
data "pfptmeta_device_settings" "settings" {
  id = pfptmeta_device_settings.settings.id
}`
)

func TestAccResourceDeviceSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("device_settings", "v1/settings/device"),
		Steps: []resource.TestStep{
			{
				Config: deviceSettingsResourceStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_device_settings.settings", "id", regexp.MustCompile("^ds-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "name", "settings-name"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "description", "settings-desc"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "apply_on_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "tunnel_mode", "full"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "ztna_always_on", "true"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "proxy_always_on", "true"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "overlay_mfa_required", "true"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "overlay_mfa_refresh_period", "20"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "vpn_login_browser", "EXTERNAL"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "session_lifetime", "15"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "session_lifetime_grace", "3"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "protocol_selection_lifetime", "10"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "search_domains.0", "example1.com"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "search_domains.1", "example2.com"),
				),
			},
			{
				Config: deviceSettingsGroup + deviceSettingsResourceStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_device_settings.settings", "id", regexp.MustCompile("^ds-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "name", "settings-name1"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "description", "settings-desc1"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "apply_on_org", "false"),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_device_settings.settings", "apply_to_entities.0",
						"pfptmeta_group.ds_group", "id"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "tunnel_mode", "split"),
					resource.TestCheckNoResourceAttr("pfptmeta_device_settings.settings", "ztna_always_on"),
					resource.TestCheckNoResourceAttr("pfptmeta_device_settings.settings", "proxy_always_on"),
					resource.TestCheckNoResourceAttr("pfptmeta_device_settings.settings", "overlay_mfa_required"),
					resource.TestCheckNoResourceAttr("pfptmeta_device_settings.settings", "overlay_mfa_refresh_period"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "vpn_login_browser", "USER_DEFINED"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "session_lifetime", "14"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "session_lifetime_grace", "4"),
					resource.TestCheckResourceAttr("pfptmeta_device_settings.settings", "protocol_selection_lifetime", "12"),
					resource.TestCheckNoResourceAttr("pfptmeta_device_settings.settings", "search_domains"),
				),
			},
		},
	})
}

func TestAccDataSourceDeviceSettings(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("device_settings", "v1/settings/device"),
		Steps: []resource.TestStep{
			{
				Config: deviceSettingsResourceStep1 + deviceSettingsDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_device_settings.settings", "id", regexp.MustCompile("^ds-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "name", "settings-name"),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "description", "settings-desc"),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "apply_on_org", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "tunnel_mode", "full"),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "ztna_always_on", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "proxy_always_on", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "overlay_mfa_required", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "overlay_mfa_refresh_period", "20"),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "vpn_login_browser", "EXTERNAL"),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "session_lifetime", "15"),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "session_lifetime_grace", "3"),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "protocol_selection_lifetime", "10"),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "search_domains.0", "example1.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_device_settings.settings", "search_domains.1", "example2.com"),
				),
			},
		},
	})
}
