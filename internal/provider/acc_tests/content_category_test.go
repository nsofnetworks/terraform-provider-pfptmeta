package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	contentCategoryStep1 = `
resource "pfptmeta_content_category" "cc" {
  name                      = "cc"
  description               = "cc desc"
  confidence_level          = "HIGH"
  forbid_uncategorized_urls = true
  types                     = ["News and Media"]
  urls                      = [".ynet.co.il"]
}
`
	contentCategoryStep2 = `
resource "pfptmeta_content_category" "cc" {
  name                      = "cc1"
  description               = "cc desc 1"
  confidence_level          = "MEDIUM"
  forbid_uncategorized_urls = false
  types                     = ["News and Media", "Sports"]
  urls                      = ["192.6.6.5", ".ynet.co.il"]
}
`
	contentCategoryDataSource = `
data "pfptmeta_content_category" "cc" {
  id = pfptmeta_content_category.cc.id
}
`
)

func TestAccResourceContentCategory(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("content_category", "v1/content_categories"),
		Steps: []resource.TestStep{
			{
				Config: contentCategoryStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_content_category.cc", "id", regexp.MustCompile("^cc-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "name", "cc"),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "description", "cc desc"),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "confidence_level", "HIGH"),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "forbid_uncategorized_urls", "true"),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "types.0", "News and Media"),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "urls.0", ".ynet.co.il"),
				),
			},
			{
				Config: contentCategoryStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_content_category.cc", "id", regexp.MustCompile("^cc-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "name", "cc1"),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "description", "cc desc 1"),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "confidence_level", "MEDIUM"),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "forbid_uncategorized_urls", "false"),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "types.0", "News and Media"),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "types.1", "Sports"),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "urls.0", "192.6.6.5"),
					resource.TestCheckResourceAttr("pfptmeta_content_category.cc", "urls.1", ".ynet.co.il"),
				),
			},
		},
	})
}

func TestAccDataSourceContentCategory(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: contentCategoryStep1 + contentCategoryDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_content_category.cc", "id", regexp.MustCompile("^cc-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_content_category.cc", "name", "cc"),
					resource.TestCheckResourceAttr("data.pfptmeta_content_category.cc", "description", "cc desc"),
					resource.TestCheckResourceAttr("data.pfptmeta_content_category.cc", "confidence_level", "HIGH"),
					resource.TestCheckResourceAttr("data.pfptmeta_content_category.cc", "forbid_uncategorized_urls", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_content_category.cc", "types.0", "News and Media"),
					resource.TestCheckResourceAttr("data.pfptmeta_content_category.cc", "urls.0", ".ynet.co.il"),
				),
			},
		},
	})
}
