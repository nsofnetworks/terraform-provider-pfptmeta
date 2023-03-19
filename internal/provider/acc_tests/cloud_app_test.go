package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	cloudAppStep1 = `
resource "pfptmeta_cloud_app" "ca" {
  name        = "cloud app"
  description = "cloud app description"
  app         = "sia-xVb7vgt6rKoQnY4Rj"
  urls        = ["ynet.co.il"]
}
`
	cloudAppStep2 = `
resource "pfptmeta_cloud_app" "ca" {
  name        = "cloud app 1"
  description = "cloud app description 1"
  app         = "sia-K8jJPltjdMv39A4wx"
  urls        = ["192.6.6.5", "ynet.co.il"]
}
`
	cloudAppDataSource = `
data "pfptmeta_cloud_app" "ca" {
  id = pfptmeta_cloud_app.ca.id
}
`
)

func TestAccResourceCloudApp(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("cloud_app", "v1/cloud_apps"),
		Steps: []resource.TestStep{
			{
				Config: cloudAppStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_cloud_app.ca", "id", regexp.MustCompile("^ca-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_cloud_app.ca", "name", "cloud app"),
					resource.TestCheckResourceAttr("pfptmeta_cloud_app.ca", "description", "cloud app description"),
					resource.TestCheckResourceAttr("pfptmeta_cloud_app.ca", "app", "sia-xVb7vgt6rKoQnY4Rj"),
					resource.TestCheckResourceAttr("pfptmeta_cloud_app.ca", "urls.0", "ynet.co.il"),
				),
			},
			{
				Config: cloudAppStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_cloud_app.ca", "id", regexp.MustCompile("^ca-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_cloud_app.ca", "name", "cloud app 1"),
					resource.TestCheckResourceAttr("pfptmeta_cloud_app.ca", "description", "cloud app description 1"),
					resource.TestCheckResourceAttr("pfptmeta_cloud_app.ca", "app", "sia-K8jJPltjdMv39A4wx"),
					resource.TestCheckResourceAttr("pfptmeta_cloud_app.ca", "urls.0", "192.6.6.5"),
					resource.TestCheckResourceAttr("pfptmeta_cloud_app.ca", "urls.1", "ynet.co.il"),
				),
			},
		},
	})
}

func TestAccDataSourceCloudApp(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: cloudAppStep1 + cloudAppDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_cloud_app.ca", "id", regexp.MustCompile("^ca-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_cloud_app.ca", "name", "cloud app"),
					resource.TestCheckResourceAttr("pfptmeta_cloud_app.ca", "description", "cloud app description"),
					resource.TestCheckResourceAttr("pfptmeta_cloud_app.ca", "app", "sia-xVb7vgt6rKoQnY4Rj"),
					resource.TestCheckResourceAttr("pfptmeta_cloud_app.ca", "urls.0", "ynet.co.il"),
				),
			},
		},
	})
}
