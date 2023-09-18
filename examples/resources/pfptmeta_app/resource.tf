variable "mappedAttributes" {
  type = list(map(string))
  default = [
    {
      attribute_format     = "unspecified"
      target_variable_name = "first_name"
      variable_name        = "given_name"
    },
    {
      attribute_format     = "unspecified"
      target_variable_name = "last_name"
      variable_name        = "family_name"
    },
    {
      target_variable_name = "groups"
      variable_name        = "groups"
      filter_type          = "equals"
      filter_value         = "grp1"
    },
    {
      variable_name = "tags"
      filter_type   = "all"
    },
  ]
}

resource "pfptmeta_app" "app_saml" {
  name             = "saml app name"
  description      = "saml app description"
  enabled          = true
  visible          = true
  ip_whitelist     = ["1.1.1.1"]
  direct_sso_login = "idp-abcd1234"
  assigned_members = ["usr-abcd1234"]
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

resource "pfptmeta_app" "app_oidc" {
  name             = "oidc app name"
  description      = "oidc app description"
  enabled          = true
  visible          = true
  ip_whitelist     = ["1.1.1.1"]
  assigned_members = ["usr-abcd1234"]
  protocol         = "OIDC"

  oidc {
    sign_in_redirect_urls = ["https://login.myApp.example.com"]
    grant_types           = ["authorization_code"]
    scopes                = ["openid", "profile", "email"]
    initiate_login_url    = "https://intial-login.myApp.example.com"
  }
}