package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	tenantRestrictionGoogleStep1 = `
resource "pfptmeta_tenant_restriction" "google" {
  name        = "google tr"
  description = "google tr desc"
  google_config {
    allow_consumer_access  = true
    allow_service_accounts = true
    tenants                = ["altostrat.com"]
  }
}
`
	tenantRestrictionMicrosoftStep1 = `
resource "pfptmeta_tenant_restriction" "microsoft" {
  name        = "microsoft tr"
  description = "microsoft tr desc"
  microsoft_config {
    allow_personal_microsoft_domains = true
    tenant_directory_id              = "456ff232-35l2-5h23-b3b3-3236w0826f3d"
    tenants                          = ["onmicrosoft.com"]
  }
}
`
	tenantRestrictionGoogleStep2 = `
resource "pfptmeta_tenant_restriction" "google" {
  name        = "google tr 1"
  description = "google tr desc 1"
  google_config {
    allow_service_accounts = false
    allow_consumer_access = false
    tenants                = ["tenorstrat.com", "altostrat.com"]
  }
}
`
	tenantRestrictionMicrosoftStep2 = `
resource "pfptmeta_tenant_restriction" "microsoft" {
  name        = "microsoft tr 1"
  description = "microsoft tr desc 1"
  microsoft_config {
    allow_personal_microsoft_domains = false
    tenant_directory_id              = "72f988bf-86f1-41af-91ab-2d7cd011db47"
    tenants                          = ["456ff232-35l2-5h23-b3b3-3236w0826f3d"]
  }
}
`
	tenantRestrictionDataSource = `
data "pfptmeta_tenant_restriction" "google" {
  id = pfptmeta_tenant_restriction.google.id
}
`
)

func TestAccResourceTenantRestriction(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("tenant_restriction", "v1/tenant_restrictions"),
		Steps: []resource.TestStep{
			{
				Config: tenantRestrictionGoogleStep1 + tenantRestrictionMicrosoftStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_tenant_restriction.google", "id", regexp.MustCompile("^tr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "name", "google tr"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "description", "google tr desc"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "google_config.0.allow_consumer_access", "true"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "google_config.0.allow_service_accounts", "true"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "google_config.0.tenants.0", "altostrat.com"),

					resource.TestMatchResourceAttr("pfptmeta_tenant_restriction.microsoft", "id", regexp.MustCompile("^tr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.microsoft", "name", "microsoft tr"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.microsoft", "description", "microsoft tr desc"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.microsoft",
						"microsoft_config.0.allow_personal_microsoft_domains", "true"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.microsoft",
						"microsoft_config.0.tenant_directory_id", "456ff232-35l2-5h23-b3b3-3236w0826f3d"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.microsoft",
						"microsoft_config.0.tenants.0", "onmicrosoft.com"),
				),
			},
			{
				Config: tenantRestrictionGoogleStep2 + tenantRestrictionMicrosoftStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_tenant_restriction.google", "id", regexp.MustCompile("^tr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "name", "google tr 1"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "description", "google tr desc 1"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "google_config.0.allow_consumer_access", "false"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "google_config.0.allow_service_accounts", "false"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "google_config.0.tenants.0", "tenorstrat.com"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "google_config.0.tenants.1", "altostrat.com"),

					resource.TestMatchResourceAttr("pfptmeta_tenant_restriction.microsoft", "id", regexp.MustCompile("^tr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.microsoft", "name", "microsoft tr 1"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.microsoft", "description", "microsoft tr desc 1"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.microsoft",
						"microsoft_config.0.allow_personal_microsoft_domains", "false"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.microsoft",
						"microsoft_config.0.tenant_directory_id", "72f988bf-86f1-41af-91ab-2d7cd011db47"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.microsoft",
						"microsoft_config.0.tenants.0", "456ff232-35l2-5h23-b3b3-3236w0826f3d"),
				),
			},
		},
	})
}

func TestAccDataSourceTenantRestriction(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("tenant_restriction", "v1/tenant_restrictions"),
		Steps: []resource.TestStep{
			{
				Config: tenantRestrictionGoogleStep1 + tenantRestrictionDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_tenant_restriction.google", "id", regexp.MustCompile("^tr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "name", "google tr"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "description", "google tr desc"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "google_config.0.allow_consumer_access", "true"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "google_config.0.allow_service_accounts", "true"),
					resource.TestCheckResourceAttr("pfptmeta_tenant_restriction.google", "google_config.0.tenants.0", "altostrat.com"),
				),
			},
		},
	})
}
