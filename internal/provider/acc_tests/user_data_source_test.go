package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

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
