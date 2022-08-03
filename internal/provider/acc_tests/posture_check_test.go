package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	userConf = `
resource "pfptmeta_user" "user" {
  given_name  = "John"
  family_name = "Smith"
  email       = "john.smith@example.com"
}
`
	postureCheckStep1 = `
resource "pfptmeta_posture_check" "check" {
  name                 = "check-name"
  description          = "check-desc"
  apply_to_entities    = [pfptmeta_user.user.id]
  osquery              = "select * from processes where name='falcon-sensor' and state='S';"
  platform             = "Linux"
  enabled              = true
  action               = "WARNING"
  when                 = ["PERIODIC"]
  interval             = 60
  user_message_on_fail = "check failed"
}
`
	postureCheckStep2 = `
resource "pfptmeta_posture_check" "check" {
  name                 = "check-name1"
  apply_to_org         = true
  exempt_entities      = [pfptmeta_user.user.id]
  check {
    type        = "minimum_app_version"
    min_version = "4.0.0"
  }
  platform             = "iOS"
  action               = "DISCONNECT"
  when                 = ["PRE_CONNECT"]
  user_message_on_fail = "check failed1"
}
`
	postureCheckStep3 = `
resource "pfptmeta_posture_check" "check" {
  name                 = "check-name2"
  apply_to_org         = true
  exempt_entities      = [pfptmeta_user.user.id]
  check {
    type        = "minimum_app_version"
    min_version = "4.0.0"
  }
  platform             = "iOS"
  action               = "WARNING"
  when                 = ["PERIODIC"]
  interval             = 5
  user_message_on_fail = "check failed2"
}
`
	dataSourcePostureCheck = `
data "pfptmeta_posture_check" "check" {
  id = pfptmeta_posture_check.check.id
}
`
)

func TestAccResourcePostureCheck(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("posture_check", "v1/posture_checks"),
		Steps: []resource.TestStep{
			{
				Config: userConf + postureCheckStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_posture_check.check", "id", regexp.MustCompile("^pc-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "name", "check-name"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "description", "check-desc"),
					resource.TestMatchResourceAttr("pfptmeta_posture_check.check", "apply_to_entities.0", regexp.MustCompile("^usr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "osquery", "select * from processes where name='falcon-sensor' and state='S';"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "platform", "Linux"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "enabled", "true"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "action", "WARNING"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "when.0", "PERIODIC"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "interval", "60"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "user_message_on_fail", "check failed"),
				),
			},
			{
				Config: userConf + postureCheckStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_posture_check.check", "id", regexp.MustCompile("^pc-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "name", "check-name1"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "apply_to_org", "true"),
					resource.TestMatchResourceAttr("pfptmeta_posture_check.check", "exempt_entities.0", regexp.MustCompile("^usr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "platform", "iOS"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "enabled", "true"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "action", "DISCONNECT"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "when.0", "PRE_CONNECT"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "user_message_on_fail", "check failed1"),
				),
			},
			{
				Config: userConf + postureCheckStep3,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_posture_check.check", "id", regexp.MustCompile("^pc-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "name", "check-name2"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "apply_to_org", "true"),
					resource.TestMatchResourceAttr("pfptmeta_posture_check.check", "exempt_entities.0", regexp.MustCompile("^usr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "platform", "iOS"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "enabled", "true"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "action", "WARNING"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "when.0", "PERIODIC"),
					resource.TestCheckResourceAttr("pfptmeta_posture_check.check", "user_message_on_fail", "check failed2"),
				),
			},
		},
	})
}

func TestAccDataSourcePostureCheck(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("posture_check", "v1/posture_checks"),
		Steps: []resource.TestStep{
			{
				Config: userConf + postureCheckStep1 + dataSourcePostureCheck,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_posture_check.check", "id", regexp.MustCompile("^pc-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_posture_check.check", "name", "check-name"),
					resource.TestCheckResourceAttr("data.pfptmeta_posture_check.check", "description", "check-desc"),
					resource.TestMatchResourceAttr("data.pfptmeta_posture_check.check", "apply_to_entities.0", regexp.MustCompile("^usr-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_posture_check.check", "osquery", "select * from processes where name='falcon-sensor' and state='S';"),
					resource.TestCheckResourceAttr("data.pfptmeta_posture_check.check", "platform", "Linux"),
					resource.TestCheckResourceAttr("data.pfptmeta_posture_check.check", "enabled", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_posture_check.check", "action", "WARNING"),
					resource.TestCheckResourceAttr("data.pfptmeta_posture_check.check", "when.0", "PERIODIC"),
					resource.TestCheckResourceAttr("data.pfptmeta_posture_check.check", "interval", "60"),
					resource.TestCheckResourceAttr("data.pfptmeta_posture_check.check", "user_message_on_fail", "check failed"),
				),
			},
		},
	})
}
