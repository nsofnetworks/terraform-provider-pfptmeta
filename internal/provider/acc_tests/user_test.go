package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccResourceUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("user", "v1/users"),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_user.user", "id", regexp.MustCompile("^usr-.*$")),
					resource.TestCheckResourceAttr("pfptmeta_user.user", "name", "John Smith"),
					resource.TestCheckResourceAttr("pfptmeta_user.user", "email", "john.smith@example.com"),
					resource.TestCheckResourceAttr("pfptmeta_user.user", "description", "u-description"),
					resource.TestCheckResourceAttr("pfptmeta_user.user", "phone", "+97251234567"),
					resource.TestCheckResourceAttr("pfptmeta_user.user", "tags.tag_name1", "tag_value1"),
					resource.TestCheckResourceAttr("pfptmeta_user.user", "tags.tag_name2", "tag_value2"),
				),
			},
			{
				Config: testAccResourceUserStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_user.user", "id", regexp.MustCompile("^usr-.*$")),
					resource.TestCheckResourceAttr("pfptmeta_user.user", "name", "John1 Smith1"),
					resource.TestCheckResourceAttr("pfptmeta_user.user", "email", "john1.smith1@example.com"),
					resource.TestCheckResourceAttr("pfptmeta_user.user", "description", "u-description1"),
					resource.TestCheckResourceAttr("pfptmeta_user.user", "phone", ""),
					resource.TestCheckResourceAttr("pfptmeta_user.user", "tags.tag_name1", "tag_value11"),
					resource.TestCheckResourceAttr("pfptmeta_user.user", "tags.tag_name2", "tag_value22"),
				),
			},
		},
	})
}

const testAccResourceUserStep1 = `
resource "pfptmeta_user" "user" {
  given_name  = "John"
  family_name = "Smith"
  email       = "john.smith@example.com"
  description = "u-description"
  phone       = "+97251234567"
  tags        = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
}
`

const testAccResourceUserStep2 = `
resource "pfptmeta_user" "user" {
  given_name  = "John1"
  family_name = "Smith1"
  email       = "john1.smith1@example.com"
  description = "u-description1"
  tags        = {
    tag_name1 = "tag_value11"
    tag_name2 = "tag_value22"
  }
}
`

func TestAccDataSourceUser(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceUser,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.pfptmeta_user.user_by_id", "id", "usr-LdjvfnK5713B8K1",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_user.user_by_id", "name", "tf-user tf-user",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_user.user_by_id", "email", "tf-user@proofpoint.com",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_user.user_by_id", "given_name", "tf-user",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_user.user_by_id", "family_name", "tf-user",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_user.user_by_email", "id", "usr-LdjvfnK5713B8K1",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_user.user_by_email", "name", "tf-user tf-user",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_user.user_by_email", "email", "tf-user@proofpoint.com",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_user.user_by_id", "given_name", "tf-user",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_user.user_by_id", "family_name", "tf-user",
					),
				),
			},
		},
	})
}

const testAccDataSourceUser = `
data "pfptmeta_user" "user_by_id" {
  id = "usr-LdjvfnK5713B8K1"
}

data "pfptmeta_user" "user_by_email" {
  email = "tf-user@proofpoint.com"
}
`
