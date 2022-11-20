package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccResourceGroupRolesAttachment(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("group", "v1/groups"),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupRoleAttachmentStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_group_roles_attachment.attachment", "id", regexp.MustCompile("^grp-.*$")),
					resource.TestMatchResourceAttr("pfptmeta_group_roles_attachment.attachment", "group_id", regexp.MustCompile("^grp-.*$")),
					resource.TestMatchResourceAttr("pfptmeta_group_roles_attachment.attachment", "roles.0", regexp.MustCompile("^rol-.*$")),
					resource.TestMatchResourceAttr("pfptmeta_group_roles_attachment.attachment", "roles.1", regexp.MustCompile("^rol-.*$")),
				),
			},
			{
				Config: testAccResourceGroupRoleAttachmentStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_group_roles_attachment.attachment", "id", regexp.MustCompile("^grp-.*$")),
					resource.TestMatchResourceAttr("pfptmeta_group_roles_attachment.attachment", "group_id", regexp.MustCompile("^grp-.*$")),
					resource.TestMatchResourceAttr("pfptmeta_group_roles_attachment.attachment", "roles.0", regexp.MustCompile("^rol-.*$")),
					resource.TestCheckNoResourceAttr("pfptmeta_group_roles_attachment.attachment", "roles.1"),
				),
			},
		},
	})
}

const testAccResourceGroupRoleAttachmentStep1 = `
resource "pfptmeta_group" "group" {
  name = "admins"
}

resource "pfptmeta_role" "metaport_role" {
  name       = "metaport role"
  apply_to_orgs = ["org-31126"]
  privileges = ["metaports:read", "metaports:write"]
}

resource "pfptmeta_role" "network_element_role" {
  name       = "network element role"
  apply_to_orgs = ["org-31126"]
  privileges = ["network_elements:read", "network_elements:write"]
}

resource "pfptmeta_group_roles_attachment" "attachment" {
  group_id = pfptmeta_group.group.id
  roles    = [pfptmeta_role.metaport_role.id, pfptmeta_role.network_element_role.id]
}
`

const testAccResourceGroupRoleAttachmentStep2 = `
resource "pfptmeta_group" "group" {
  name = "admins"
}

resource "pfptmeta_role" "metaport_role" {
  name       = "metaport role"
  apply_to_orgs = ["org-31126"]
  privileges = ["metaports:read", "metaports:write"]
}

resource "pfptmeta_role" "network_element_role" {
  name       = "network element role"
  apply_to_orgs = ["org-31126"]
  privileges = ["network_elements:read", "network_elements:write"]
}

resource "pfptmeta_group_roles_attachment" "attachment" {
  group_id = pfptmeta_group.group.id
  roles    = [pfptmeta_role.metaport_role.id]
}
`
