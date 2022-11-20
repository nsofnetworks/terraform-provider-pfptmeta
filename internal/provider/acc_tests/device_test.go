package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceDeviceEntity(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("device", "v1/devices"),
		Steps: []resource.TestStep{
			{
				Config: TestAccResourceDeviceConfig + deviceDependencies,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_device.device", "id", regexp.MustCompile("^dev-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_device.device", "name", "resourceDeviceTest",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_device.device", "description", "some details about the device",
					),
					resource.TestMatchResourceAttr(
						"pfptmeta_device.device", "owner_id", regexp.MustCompile("^usr-.*$"),
					),
				),
			},
		},
	})
}

func TestAccDataSourceDeviceEntity(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{Config: TestAccResourceDeviceConfig + deviceDependencies},
			{
				Config: TestAccResourceDeviceConfig + deviceDataSource + deviceDependencies,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_device.device", "id", regexp.MustCompile("^dev-.*$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_device.device", "name", "resourceDeviceTest",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_device.device", "description", "some details about the device",
					),
					resource.TestMatchResourceAttr(
						"data.pfptmeta_device.device", "owner_id", regexp.MustCompile("^usr-.*$"),
					),
				),
			},
		},
	})
}

const (
	deviceDependencies = `
resource "pfptmeta_user" "user" {
	given_name  = "John-dev"
	family_name = "Doe-dev"
	email       = "john.doe@example.com"
}
`
	TestAccResourceDeviceConfig = `
resource "pfptmeta_device" "device" {
	name        = "resourceDeviceTest"
	description = "some details about the device"
	owner_id    =  pfptmeta_user.user.id
}
`
	deviceDataSource = `
data "pfptmeta_device" "device" {
	id = pfptmeta_device.device.id
}
`
)
