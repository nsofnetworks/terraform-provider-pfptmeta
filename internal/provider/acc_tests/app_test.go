package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const testAppDependencies = `
data "pfptmeta_user" "app_user_by_email" {
  email = "tf-user@proofpoint.com"
}
`

const testAppSamlResource = `
resource "pfptmeta_app" "app_saml" {
  name             = "saml app name"
  description      = "saml app description"
  enabled          = true
  visible          = true
  ip_whitelist     = ["1.1.1.1"]
  assigned_members = [data.pfptmeta_user.app_user_by_email.id]
  protocol         = "SAML"

  saml {
    audience_uri              = "https://audience.myApp.com"
    sso_acs_url               = "https://login.myApp.example.com"
    destination               = "https://login.myApp.example.com"
    recipient                 = "https://login.myApp.example.com"
    default_relay_state       = "https://relay.myApp.com"
    subject_name_id_attribute = "email"
    subject_name_id_format    = "emailAddress"
    signature_algorithm       = "RSA-SHA256"
    digest_algorithm          = "SHA256"
  }
}

data "pfptmeta_app" "app_saml" {
  id       = pfptmeta_app.app_saml.id
  protocol = "SAML"
}
`

const testAppOidcResource = `
variable "mappedAttributes" {
  type = list(map(string))
  default = [
    {
      variable_name = "tags"
      filter_type   = "all"
    },
  ]
}

resource "pfptmeta_app" "app_oidc" {
  name             = "oidc app name"
  description      = "oidc app description"
  enabled          = true
  visible          = true
  ip_whitelist     = ["1.1.1.1"]
  assigned_members = [data.pfptmeta_user.app_user_by_email.id]
  protocol         = "OIDC"

  oidc {
    sign_in_redirect_urls = ["https://login.myApp.example.com"]
    grant_types           = ["authorization_code"]
    scopes                = ["openid"]
    initiate_login_url    = "https://intial-login.myApp.example.com"
  }

  dynamic "mapped_attributes" {
    for_each = var.mappedAttributes
    content {
      variable_name        = mapped_attributes.value.variable_name
      attribute_format     = try(mapped_attributes.value.attribute_format, null)
      target_variable_name = try(mapped_attributes.value.target_variable_name, null)
      filter_type          = try(mapped_attributes.value.filter_type, null)
      filter_value         = try(mapped_attributes.value.filter_value, null)
    }
  }
}

data "pfptmeta_app" "app_oidc" {
  id       = pfptmeta_app.app_oidc.id
  protocol = "OIDC"
}
`

func TestAccDataSourceAppOidc(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("app", "v1/apps"),
		Steps: []resource.TestStep{
			{
				Config: testAppDependencies + testAppOidcResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_app.app_oidc", "id", regexp.MustCompile("^app-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "name", "oidc app name"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "description", "oidc app description"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "enabled", "true"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "visible", "true"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "ip_whitelist.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "assigned_members.0", "usr-xN6MCvzmWyvJYdk"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "oidc.0.sign_in_redirect_urls.0",
						"https://login.myApp.example.com"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "oidc.0.grant_types.0", "authorization_code"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "oidc.0.scopes.0", "openid"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "oidc.0.initiate_login_url",
						"https://intial-login.myApp.example.com"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "oidc.0.id_token_lifetime", "30"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "oidc.0.access_token_lifetime", "5"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "mapped_attributes.0.variable_name", "tags"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "mapped_attributes.0.target_variable_name", ""),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "mapped_attributes.0.attribute_format", "unspecified"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "mapped_attributes.0.filter_type", "all"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_oidc", "mapped_attributes.0.filter_value", ""),
				),
			},
			{
				Config: testAppDependencies + testAppOidcResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_app.app_oidc", "id", regexp.MustCompile("^app-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "name", "oidc app name"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "description", "oidc app description"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "enabled", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "visible", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "ip_whitelist.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "assigned_members.0", "usr-xN6MCvzmWyvJYdk"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "oidc.0.sign_in_redirect_urls.0",
						"https://login.myApp.example.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "oidc.0.grant_types.0", "authorization_code"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "oidc.0.scopes.0", "openid"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "oidc.0.initiate_login_url",
						"https://intial-login.myApp.example.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "oidc.0.id_token_lifetime", "30"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "oidc.0.access_token_lifetime", "5"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "mapped_attributes.0.variable_name", "tags"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "mapped_attributes.0.target_variable_name", ""),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "mapped_attributes.0.attribute_format", "unspecified"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "mapped_attributes.0.filter_type", "all"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_oidc", "mapped_attributes.0.filter_value", ""),
				),
			},
		},
	})
}

func TestAccDataSourceAppSaml(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("app", "v1/apps"),
		Steps: []resource.TestStep{
			{
				Config: testAppDependencies + testAppSamlResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_app.app_saml", "id", regexp.MustCompile("^app-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "name", "saml app name"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "description", "saml app description"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "enabled", "true"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "visible", "true"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "ip_whitelist.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "assigned_members.0", "usr-xN6MCvzmWyvJYdk"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "saml.0.audience_uri", "https://audience.myApp.com"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "saml.0.sso_acs_url", "https://login.myApp.example.com"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "saml.0.destination", "https://login.myApp.example.com"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "saml.0.default_relay_state", "https://relay.myApp.com"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "saml.0.subject_name_id_attribute", "email"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "saml.0.subject_name_id_format", "emailAddress"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "saml.0.signature_algorithm", "RSA-SHA256"),
					resource.TestCheckResourceAttr("pfptmeta_app.app_saml", "saml.0.digest_algorithm", "SHA256"),
				),
			},
			{
				Config: testAppDependencies + testAppSamlResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_app.app_saml", "id", regexp.MustCompile("^app-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "name", "saml app name"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "description", "saml app description"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "enabled", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "visible", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "ip_whitelist.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "assigned_members.0", "usr-xN6MCvzmWyvJYdk"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "saml.0.audience_uri", "https://audience.myApp.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "saml.0.sso_acs_url", "https://login.myApp.example.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "saml.0.destination", "https://login.myApp.example.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "saml.0.default_relay_state", "https://relay.myApp.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "saml.0.subject_name_id_attribute", "email"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "saml.0.subject_name_id_format", "emailAddress"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "saml.0.signature_algorithm", "RSA-SHA256"),
					resource.TestCheckResourceAttr("data.pfptmeta_app.app_saml", "saml.0.digest_algorithm", "SHA256"),
				),
			},
		},
	})
}
