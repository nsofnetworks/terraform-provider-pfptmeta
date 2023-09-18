package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const testIdpSamlResource = `
variable "saml_cert" {
  type    = string
  default = <<EOF
-----BEGIN CERTIFICATE-----
MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbty
WPT9pURI3lbsKMiB6Fn/VHOKE13p4D8xgOCADpdRagdT6n4etr9atzDKUSvpMtR3
CP5noNc97WiNCggBjVWhs7szEe8ugyqF23XwpHQ6uV1LKH50m92MbOWfCtjU9p/x
qhNpQQ1AZhqNy5Gevap5k8XzRmjSldNAFZMY7Yv3Gi+nyCwGwpVtBUwhuLzgNFK/
yDtw2WcWmUU7NuC8Q6MWvPebxVtCfVp/iQU6q60yyt6aGOBkhAX0LpKAEhKidixY
nP9PNVBvxgu3XZ4P36gZV6+ummKdBVnc3NqwBLu5+CcdRdusmHPHd5pHf4/38Z3/
6qU2a/fPvWzceVTEgZ47QjFMTCTmCwNt29cvi7zZeQzjtwQgn4ipN9NibRH/Ax/q
TbIzHfrJ1xa2RteWSdFjwtxi9C20HUkjXSeI4YlzQMH0fPX6KCE7aVePTOnB69I/
a9/q96DiXZajwlpq3wFctrs1oXqBp5DVrCIj8hU2wNgB7LtQ1mCtsYz//heai0K9
PhE4X6hiE0YmeAZjR0uHl8M/5aW9xCoJ72+12kKpWAa0SFRWLy6FejNYCYpkupVJ
yecLk/4L1W0l6jQQZnWErXZYe0PNFcmwGXy1Rep83kfBRNKRy5tvocalLlwXLdUk
AIU+2GKjyT3iMuzZxxFxPFMCAwEAAQ==
-----END CERTIFICATE-----
EOF
}

resource "pfptmeta_idp" "idp_saml" {
  name              = "saml idp name"
  description       = "saml idp description"
  enabled           = true
  hidden            = false
  icon              = "my_icon.svg"
  mapped_attributes = ["immutable_id"]

  saml_config {
    issuer              = "https://issuer.myIdp.com"
    certificate         = var.saml_cert
    sso_url             = "https://myIdp.login.com"
    authn_context_class = "PasswordProtectedTransport"
    jit_enabled         = false
  }

  scim_config {
    api_key_id       = "key-rNPQCZodWRwNb73"
    assume_ownership = false
  }
}

data "pfptmeta_idp" "idp_saml" {
  id       = pfptmeta_idp.idp_saml.id
}
`

const testIdpScimResource = `
resource "pfptmeta_idp" "idp_scim" {
  name        = "scim idp name"
  description = "scim idp description"
  enabled     = true
  hidden      = true

  scim_config {
    api_key_id       = "key-rNPQCZodWRwNb73"
    assume_ownership = true
  }
}

data "pfptmeta_idp" "idp_scim" {
  id       = pfptmeta_idp.idp_scim.id
}
`

func TestAccDataSourceIdpSaml(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("idp", "v1/settings/idps"),
		Steps: []resource.TestStep{
			{
				Config: testIdpSamlResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_idp.idp_saml", "id", regexp.MustCompile("^idp-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_saml", "name", "saml idp name"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_saml", "description", "saml idp description"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_saml", "enabled", "true"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_saml", "hidden", "false"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_saml", "mapped_attributes.0", "immutable_id"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_saml", "saml_config.0.issuer", "https://issuer.myIdp.com"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_saml", "saml_config.0.sso_url", "https://myIdp.login.com"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_saml", "saml_config.0.authn_context_class", "PasswordProtectedTransport"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_saml", "saml_config.0.jit_enabled", "false"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_saml", "scim_config.0.api_key_id", "key-rNPQCZodWRwNb73"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_saml", "scim_config.0.assume_ownership", "false"),
				),
			},
			{
				Config: testIdpSamlResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_idp.idp_saml", "id", regexp.MustCompile("^idp-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_saml", "name", "saml idp name"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_saml", "description", "saml idp description"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_saml", "enabled", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_saml", "hidden", "false"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_saml", "mapped_attributes.0", "immutable_id"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_saml", "saml_config.0.issuer", "https://issuer.myIdp.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_saml", "saml_config.0.sso_url", "https://myIdp.login.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_saml", "saml_config.0.authn_context_class", "PasswordProtectedTransport"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_saml", "saml_config.0.jit_enabled", "false"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_saml", "scim_config.0.api_key_id", "key-rNPQCZodWRwNb73"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_saml", "scim_config.0.assume_ownership", "false"),
				),
			},
		},
	})
}

func TestAccDataSourceIdpScim(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("idp", "v1/settings/idps"),
		Steps: []resource.TestStep{
			{
				Config: testIdpScimResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_idp.idp_scim", "id", regexp.MustCompile("^idp-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_scim", "name", "scim idp name"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_scim", "description", "scim idp description"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_scim", "enabled", "true"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_scim", "hidden", "true"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_scim", "scim_config.0.api_key_id", "key-rNPQCZodWRwNb73"),
					resource.TestCheckResourceAttr("pfptmeta_idp.idp_scim", "scim_config.0.assume_ownership", "true"),
				),
			},
			{
				Config: testIdpScimResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_idp.idp_scim", "id", regexp.MustCompile("^idp-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_scim", "name", "scim idp name"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_scim", "description", "scim idp description"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_scim", "enabled", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_scim", "hidden", "true"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_scim", "scim_config.0.api_key_id", "key-rNPQCZodWRwNb73"),
					resource.TestCheckResourceAttr("data.pfptmeta_idp.idp_scim", "scim_config.0.assume_ownership", "true"),
				),
			},
		},
	})
}
