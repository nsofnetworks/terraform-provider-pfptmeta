package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDeviceAlias(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDeviceAliasStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_device_alias.alias", "device_id", regexp.MustCompile("^dev-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_device_alias.alias", "alias", "alias.test.com",
					),
				),
			},
		},
	})
}

func TestAccDataSourceDeviceAlias(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceDeviceAliasStep1 + testAccResourceDeviceAliasDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_device_alias.alias", "device_id", regexp.MustCompile("^dev-.*$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_device_alias.alias", "alias", "alias.test.com",
					),
				),
			},
		},
	})
}

const testAccResourceDeviceAliasStep1 = `
resource "pfptmeta_user" "user" {
	given_name  = "Abc"
	family_name = "Def"
	email       = "abc.def@example.com"
}

resource "pfptmeta_device" "device" {
  name        = "dummy device"
  description = "some details about the device"
  owner_id    =  pfptmeta_user.user.id
}

resource "pfptmeta_device_alias" "alias" {
  device_id = pfptmeta_device.device.id
  alias     = "alias.test.com"
}
`

const testAccResourceDeviceAliasDataSource = `
data "pfptmeta_device_alias" "alias" {
  device_id = pfptmeta_device_alias.alias.device_id
  alias     = "alias.test.com"
}
`
