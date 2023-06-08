package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const catalogAppDataSource = `
data "pfptmeta_catalog_app" "catalog_app" {
  name     = "Google Earth"
  category = "Entertainment and Lifestyle"
}
`

func TestAccDataSourceCatalogApp(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: catalogAppDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.pfptmeta_catalog_app.catalog_app", "id", "sia-GG8JA7fDq4PYxYQ4q"),
					resource.TestCheckResourceAttr("data.pfptmeta_catalog_app.catalog_app", "name", "Google Earth"),
					resource.TestCheckResourceAttr("data.pfptmeta_catalog_app.catalog_app", "risk", "5"),
					resource.TestCheckResourceAttr("data.pfptmeta_catalog_app.catalog_app", "urls.0", "earth.google.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_catalog_app.catalog_app", "vendor", "Google, LLC"),
					resource.TestCheckResourceAttr("data.pfptmeta_catalog_app.catalog_app", "verified", "false"),
				),
			},
		},
	})
}
