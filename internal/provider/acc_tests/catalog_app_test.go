package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

const catalogAppDataSource = `
data "pfptmeta_catalog_app" "catalog_app" {
  name     = "Google"
  category = "Property Management"
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
					resource.TestCheckResourceAttr("data.pfptmeta_content_category.catalog_app", "id", "sia-8QB18xfy75jknrPk4"),
					resource.TestCheckResourceAttr("data.pfptmeta_content_category.catalog_app", "name", "Google"),
					resource.TestCheckResourceAttr("data.pfptmeta_content_category.catalog_app", "risk", "6"),
					resource.TestCheckResourceAttr("data.pfptmeta_content_category.catalog_app", "urls.0", "sites.google.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_content_category.catalog_app", "vendor", "Google, LLC"),
					resource.TestCheckResourceAttr("data.pfptmeta_content_category.catalog_app", "verified", "false"),
				),
			},
		},
	})
}
