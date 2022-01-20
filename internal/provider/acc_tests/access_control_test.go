package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	accessControlDependencies = `
resource "pfptmeta_group" "applied_group" {
  name = "access-control-group"
}
resource "pfptmeta_group" "exempt_group" {
  name = "exempt-access-control-group"
}
`
	accessControlResourceStep1 = `
resource "pfptmeta_access_control" "access" {
  name            = "access control name"
  description     = "access control description"
  apply_to_org    = true
  allowed_routes  = ["192.0.2.0/24", "192.168.0.0/24"]
}
`
	accessControlResourceStep2 = `
resource "pfptmeta_access_control" "access" {
  name              = "access control name1"
  description       = "access control description1"
  apply_to_entities = [pfptmeta_group.applied_group.id]
  exempt_entities   = [pfptmeta_group.exempt_group.id]
  allowed_routes    = ["192.168.0.0/24", "192.0.2.0/24"]
}
`
	accessControlDataSource = `
data "pfptmeta_access_control" "access" {
  id = pfptmeta_access_control.access.id
}
`
)

func TestAccResourceAccessControl(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("access_control", "v1/access_controls"),
		Steps: []resource.TestStep{
			{
				Config: accessControlDependencies + accessControlResourceStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_access_control.access", "id", regexp.MustCompile("^ac-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "name", "access control name"),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "description", "access control description"),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "allowed_routes.0", "192.0.2.0/24"),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "allowed_routes.1", "192.168.0.0/24"),
				),
			},
			{
				Config: accessControlDependencies + accessControlResourceStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_access_control.access", "id", regexp.MustCompile("^ac-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "name", "access control name1"),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "description", "access control description1"),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "apply_to_org", "false"),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "allowed_routes.0", "192.168.0.0/24"),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "allowed_routes.1", "192.0.2.0/24"),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_access_control.access", "exempt_entities.0",
						"pfptmeta_group.exempt_group", "id"),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_access_control.access", "apply_to_entities.0",
						"pfptmeta_group.applied_group", "id"),
				),
			},
		},
	})
}

func TestAccDataSourceAccessControl(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("posture_check", "v1/posture_checks"),
		Steps: []resource.TestStep{
			{
				Config: accessControlResourceStep1 + accessControlDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_access_control.access", "id", regexp.MustCompile("^ac-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "name", "access control name"),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "description", "access control description"),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "apply_to_org", "true"),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "allowed_routes.0", "192.0.2.0/24"),
					resource.TestCheckResourceAttr("pfptmeta_access_control.access", "allowed_routes.1", "192.168.0.0/24"),
				),
			},
		},
	})
}
