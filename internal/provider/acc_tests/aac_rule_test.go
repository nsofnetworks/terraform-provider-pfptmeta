package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const testAacDependencies = `
data "pfptmeta_user" "aac_user_by_email" {
  email = "tf-user@proofpoint.com"
}
`

const testAacRuleResource = `
resource "pfptmeta_aac_rule" "rule" {
  name             = "aac rule name"
  description      = "arl-description"
  enabled          = true
  priority         = 1
  action           = "allow"
  apply_all_apps   = true
  sources          = [data.pfptmeta_user.aac_user_by_email.id]
  ip_reputations   = ["tor"]
}
`
const testAacRuleDataResource = `
data "pfptmeta_aac_rule" "rule" {
  id = pfptmeta_aac_rule.rule.id
}
`

func TestAccDataSourceAacRule(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("aac_rule", "v1/aac_rules"),
		Steps: []resource.TestStep{
			{
				Config: testAacDependencies + testAacRuleResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_aac_rule.rule", "id", regexp.MustCompile("^arl-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_aac_rule.rule", "name", "aac rule name"),
					resource.TestCheckResourceAttr("pfptmeta_aac_rule.rule", "description", "arl-description"),
					resource.TestCheckResourceAttr("pfptmeta_aac_rule.rule", "apply_all_apps", "true"),
					resource.TestCheckResourceAttr("pfptmeta_aac_rule.rule", "action", "allow"),
					resource.TestCheckResourceAttr("pfptmeta_aac_rule.rule", "priority", "1"),
					resource.TestCheckResourceAttr("pfptmeta_aac_rule.rule", "sources.0", "usr-xN6MCvzmWyvJYdk"),
					resource.TestCheckResourceAttr("pfptmeta_aac_rule.rule", "ip_reputations.0", "tor"),
				),
			},
			{
				Config: testAacDependencies + testAacRuleResource + testAacRuleDataResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_aac_rule.rule", "id", regexp.MustCompile("^arl-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_aac_rule.rule", "name", "aac rule name"),
					resource.TestCheckResourceAttr("data.pfptmeta_aac_rule.rule", "description", "arl-description"),
					resource.TestCheckResourceAttr("data.pfptmeta_aac_rule.rule", "apply_all_apps", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_aac_rule.rule", "action", "allow"),
					resource.TestCheckResourceAttr("data.pfptmeta_aac_rule.rule", "priority", "1"),
					resource.TestCheckResourceAttr("data.pfptmeta_aac_rule.rule", "sources.0", "usr-xN6MCvzmWyvJYdk"),
					resource.TestCheckResourceAttr("data.pfptmeta_aac_rule.rule", "ip_reputations.0", "tor"),
				),
			},
		},
	})
}
