package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccResourceGroupUsersAttachment(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceGroupUsersAttachmentStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_group_users_attachment.attachment", "id", regexp.MustCompile("^grp-.*$")),
					resource.TestMatchResourceAttr("pfptmeta_group_users_attachment.attachment", "group_id", regexp.MustCompile("^grp-.*$")),
					resource.TestMatchResourceAttr("pfptmeta_group_users_attachment.attachment", "users.0", regexp.MustCompile("^usr-.*$")),
					resource.TestMatchResourceAttr("pfptmeta_group_users_attachment.attachment", "users.1", regexp.MustCompile("^usr-.*$")),
					resource.TestMatchResourceAttr("pfptmeta_group_users_attachment.attachment", "users.2", regexp.MustCompile("^usr-.*$")),
				),
			},
			{
				Config: testAccResourceGroupUsersAttachmentStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_group_users_attachment.attachment", "id", regexp.MustCompile("^grp-.*$")),
					resource.TestMatchResourceAttr("pfptmeta_group_users_attachment.attachment", "group_id", regexp.MustCompile("^grp-.*$")),
					resource.TestMatchResourceAttr("pfptmeta_group_users_attachment.attachment", "users.0", regexp.MustCompile("^usr-.*$")),
					resource.TestCheckNoResourceAttr("pfptmeta_group_users_attachment.attachment", "users.1"),
					resource.TestCheckNoResourceAttr("pfptmeta_group_users_attachment.attachment", "users.2"),
				),
			},
		},
	})
}

const testAccResourceGroupUsersAttachmentStep1 = `
resource "pfptmeta_group" "group" {
  name = "some-group"
}

resource "pfptmeta_user" "user1" {
  given_name  = "user"
  family_name = "one"
  email       = "user1@example.com"
}

resource "pfptmeta_user" "user2" {
  given_name  = "user"
  family_name = "two"
  email       = "user2@example.com"
}

resource "pfptmeta_user" "user3" {
  given_name  = "user"
  family_name = "three"
  email       = "user3@example.com"
}

resource "pfptmeta_group_users_attachment" "attachment" {
  group_id = pfptmeta_group.group.id
  users    = [
    pfptmeta_user.user1.id,
    pfptmeta_user.user2.id,
    pfptmeta_user.user3.id
  ]
}
`

const testAccResourceGroupUsersAttachmentStep2 = `
resource "pfptmeta_group" "group" {
  name = "some-group"
}

resource "pfptmeta_user" "user1" {
  given_name  = "user"
  family_name = "one"
  email       = "user1@example.com"
}

resource "pfptmeta_user" "user2" {
  given_name  = "user"
  family_name = "two"
  email       = "user2@example.com"
}

resource "pfptmeta_user" "user3" {
  given_name  = "user"
  family_name = "three"
  email       = "user3@example.com"
}

resource "pfptmeta_group_users_attachment" "attachment" {
  group_id = pfptmeta_group.group.id
  users    = [
    pfptmeta_user.user1.id
  ]
}
`
