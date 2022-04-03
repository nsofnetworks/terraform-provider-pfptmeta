resource "pfptmeta_device_settings" "ztna_settings" {
  name           = "ZTNA settings"
  description    = "device settings description"
  apply_on_org   = true
  tunnel_mode    = "full"
  ztna_always_on = true
}

resource "pfptmeta_group" "swg_group" {
  name = "group"
}

resource "pfptmeta_device_settings" "swg_settings" {
  name              = "SWG settings"
  apply_to_entities = [pfptmeta_group.swg_group.id]
  proxy_always_on   = true
}

locals {
  identity_provider_id = "idp-123abc"
}

resource "pfptmeta_device_settings" "login_settings" {
  name                       = "Login settings"
  apply_on_org               = true
  overlay_mfa_required       = true
  overlay_mfa_refresh_period = 20
  direct_sso                 = local.identity_provider_id
  vpn_login_browser          = "EXTERNAL"
}

resource "pfptmeta_device_settings" "session_lifetime_settings" {
  name                   = "Session LifeTime settings"
  apply_on_org           = true
  session_lifetime       = 15
  session_lifetime_grace = "3"
}

resource "pfptmeta_device_settings" "advanced_settings" {
  name                        = "Advanced settings"
  apply_on_org                = true
  auto_fqdn_domain_names      = [""]
  protocol_selection_lifetime = "10"
  search_domains              = ["example1.com", "example2.com"]
}
