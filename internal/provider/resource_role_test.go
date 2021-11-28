package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
	"regexp"
	"testing"
)

func TestAccResourceRole(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("role", "v1/roles"),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceRoleStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_role.admin_role", "id", regexp.MustCompile("^rol-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "name", "admin role",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "description", "role with all privileges",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "all_read_privileges", "true",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "all_write_privileges", "true",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "all_suborgs", "false",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.with_privileges", "name", "with privs",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.with_privileges", "privileges.0", "metaports:read",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.with_privileges", "privileges.1", "metaports:write",
					),
				),
			},
			{
				Config: testAccResourceRoleStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_role.admin_role", "id", regexp.MustCompile("^rol-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "name", "admin role1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "description", "role with all privileges1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "all_read_privileges", "false",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "all_write_privileges", "false",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_role.admin_role", "all_suborgs", "false",
					),
				),
			},
		},
	})
}

const testAccResourceRoleStep1 = `
resource "pfptmeta_role" "admin_role" {
  name                 = "admin role"
  description          = "role with all privileges"
  all_read_privileges  = true
  all_write_privileges = true
}

resource "pfptmeta_role" "with_privileges" {
  name                 = "with privs"
  privileges		   = ["metaports:read", "metaports:write"]
}
`

const testAccResourceRoleStep2 = `
resource "pfptmeta_role" "admin_role" {
  name                 = "admin role1"
  description          = "role with all privileges1"
  all_read_privileges  = false
}
`

var validPrivs = []string{
	"orgs:read",
	"orgs:write",
	"users:read",
	"users:write",
	"network_elements:read",
	"network_elements:write",
	"policies:read",
	"policies:write",
	"groups:read",
	"groups:write",
	"roles:read",
	"roles:write",
	"egress_routes:read",
	"egress_routes:write",
	"locations:read",
	"audit:read",
	"metrics:read",
	"api_keys:read",
	"api_keys:write",
	"settings:read",
	"settings:write",
	"metaports:read",
	"metaports:write",
	"routing_groups:read",
	"routing_groups:write",
	"peerings:read",
	"peerings:write",
	"metaconnects:read",
	"metaconnects:write",
	"easylinks:read",
	"easylinks:write",
	"tags:read",
	"tags:write",
	"alerts:read",
	"alerts:write",
	"posture_checks:read",
	"posture_checks:write",
	"access_bridges:read",
	"access_bridges:write",
	"trusted_networks:read",
	"trusted_networks:write",
	"certificates:read",
	"certificates:write",
	"apps:read",
	"apps:write",
	"version_controls:read",
	"version_controls:write",
	"access_controls:read",
	"access_controls:write",
	"url_filtering_rules:read",
	"url_filtering_rules:write",
	"threat_categories:read",
	"threat_categories:write",
	"content_categories:read",
	"content_categories:write",
	"pac_files:read",
	"pac_files:write",
	"metaport_clusters:read",
	"metaport_clusters:write",
	"metaport_failovers:read",
	"metaport_failovers:write",
	"cloud_apps:read",
	"cloud_apps:write",
	"file_scanning_rules:read",
	"file_scanning_rules:write",
	"ssl_bypass_rules:read",
	"ssl_bypass_rules:write",
	"enterprise_dns:read",
	"enterprise_dns:write",
	"proxy_access:write",
	"dlp_rules:read",
	"dlp_rules:write",
	"tenant_restrictions:read",
	"tenant_restrictions:write",
}

var nonValidPrivs = []string{
	"abcde",
	"read",
	"write",
	"test123:read",
}

func TestPrivilegesPattern(t *testing.T) {
	for _, priv := range validPrivs {
		assert.True(t, privilegesPattern.MatchString(priv))
	}
	for _, priv := range nonValidPrivs {
		assert.False(t, privilegesPattern.MatchString(priv))
	}
}
