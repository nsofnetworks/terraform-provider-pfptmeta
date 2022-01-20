package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccResourceUserRolesAttachment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("user", "v1/users"),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceUserRoleAttachmentStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"pfptmeta_user_roles_attachment.attachment", "user_id",
						"pfptmeta_user.user", "id"),
					resource.TestMatchResourceAttr(
						"pfptmeta_user_roles_attachment.attachment", "roles.0", regexp.MustCompile("^rol-.+$")),
					resource.TestMatchResourceAttr(
						"pfptmeta_user_roles_attachment.attachment", "roles.1", regexp.MustCompile("^rol-.+$")),
				),
			},
			{
				Config: testAccResourceUserRoleAttachmentStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"pfptmeta_user_roles_attachment.attachment", "user_id",
						"pfptmeta_user.user", "id"),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_user_roles_attachment.attachment", "roles.0",
						"pfptmeta_role.metaport_role", "id"),
					resource.TestCheckNoResourceAttr("pfptmeta_user_roles_attachment.attachment", "roles.1"),
				),
			},
		},
	})
}

const testAccResourceUserRoleAttachmentStep1 = `
resource "pfptmeta_user" "user" {
  given_name  = "John"
  family_name = "Smith"
  email       = "john.smith@example.com"
}

resource "pfptmeta_role" "metaport_role" {
  name       = "metaport role"
  privileges = ["metaports:read", "metaports:write"]
}

resource "pfptmeta_role" "network_element_role" {
  name       = "network element role"
  privileges = ["network_elements:read", "network_elements:write"]
}

resource "pfptmeta_user_roles_attachment" "attachment" {
  user_id = pfptmeta_user.user.id
  roles    = [pfptmeta_role.metaport_role.id, pfptmeta_role.network_element_role.id]
}
`

const testAccResourceUserRoleAttachmentStep2 = `
resource "pfptmeta_user" "user" {
  given_name  = "John"
  family_name = "Smith"
  email       = "john.smith@example.com"
}

resource "pfptmeta_role" "metaport_role" {
  name       = "metaport role"
  privileges = ["metaports:read", "metaports:write"]
}

resource "pfptmeta_role" "network_element_role" {
  name       = "network element role"
  privileges = ["network_elements:read", "network_elements:write"]
}

resource "pfptmeta_user_roles_attachment" "attachment" {
  user_id = pfptmeta_user.user.id
  roles    = [pfptmeta_role.metaport_role.id]
}
`
