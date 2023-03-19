package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	sslBypassRuleStep1 = `
resource "pfptmeta_ssl_bypass_rule" "rule" {
  name                      = "rule name"
  description               = "rule description"
  apply_to_org              = true
  bypass_uncategorized_urls = false
  content_types             = ["Abortion"]
  domains                   = ["twitter.com"]
  priority                  = 15
}

`
	sslBypassRuleStep2 = `
resource "pfptmeta_ssl_bypass_rule" "rule" {
  name                      = "rule name 1"
  description               = "rule description 1"
  apply_to_org              = false
  bypass_uncategorized_urls = true
  content_types             = ["Abused Drugs"]
  domains                   = ["twitter1.com"]
  priority                  = 25
}

`
	sslBypassRuleDataSource = `
resource "pfptmeta_ssl_bypass_rule" "data_source_rule" {
  name                      = "data-source rule"
  description               = "rule description"
  apply_to_org              = true
  bypass_uncategorized_urls = false
  content_types             = ["Abortion"]
  domains                   = ["twitter.com"]
  priority                  = 35
}

data "pfptmeta_ssl_bypass_rule" "rule" {
  id = pfptmeta_ssl_bypass_rule.data_source_rule.id
}
`
)

func TestAccResourceSslBypassRule(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("/ssl_bypass_rule", "v1/ssl_bypass_rules"),
		Steps: []resource.TestStep{
			{
				Config: sslBypassRuleStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_ssl_bypass_rule.rule", "id", regexp.MustCompile("^sbr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "name", "rule name"),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "description", "rule description"),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "bypass_uncategorized_urls", "false"),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "content_types.0", "Abortion"),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "domains.0", "twitter.com"),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "priority", "15"),
				),
			},
			{
				Config: sslBypassRuleStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_ssl_bypass_rule.rule", "id", regexp.MustCompile("^sbr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "name", "rule name 1"),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "description", "rule description 1"),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "apply_to_org", "false"),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "bypass_uncategorized_urls", "true"),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "content_types.0", "Abused Drugs"),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "domains.0", "twitter1.com"),
					resource.TestCheckResourceAttr("pfptmeta_ssl_bypass_rule.rule", "priority", "25"),
				),
			},
		},
	})
}

func TestAccDataSourceSslBypassRule(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: sslBypassRuleDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_ssl_bypass_rule.rule", "id", regexp.MustCompile("^sbr-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_ssl_bypass_rule.rule", "name", "data-source rule"),
					resource.TestCheckResourceAttr("data.pfptmeta_ssl_bypass_rule.rule", "description", "rule description"),
					resource.TestCheckResourceAttr("data.pfptmeta_ssl_bypass_rule.rule", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_ssl_bypass_rule.rule", "bypass_uncategorized_urls", "false"),
					resource.TestCheckResourceAttr("data.pfptmeta_ssl_bypass_rule.rule", "content_types.0", "Abortion"),
					resource.TestCheckResourceAttr("data.pfptmeta_ssl_bypass_rule.rule", "domains.0", "twitter.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_ssl_bypass_rule.rule", "priority", "35"),
				),
			},
		},
	})
}
