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
}

resource "pfptmeta_idp" "idp_oidc" {
  name        = "oidc idp name"
  description = "oidc idp description"
  enabled     = true
  hidden      = false

  oidc_config {
    issuer        = "https://issuer.myIdp.com"
    client_id     = "MyIdpClientId12345"
    client_secret = "MyIdpClientSecret12345"
    jit_enabled   = false
  }
}

resource "pfptmeta_idp" "scim_oidc" {
  name        = "scim idp name"
  description = "scim idp description"
  enabled     = true
  hidden      = true

  scim_config {
    api_key_id       = "key-abcd123456"
    assume_ownership = true
  }
}