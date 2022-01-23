resource "pfptmeta_user_settings" "swg_settings" {
  name         = "SWG settings"
  description  = "device settings description"
  apply_on_org = true
  proxy_pops   = "POPS_WITH_DEDICATED_IPS"
}

resource "pfptmeta_group" "ztna_group" {
  name = "group"
}

resource "pfptmeta_user_settings" "ztna_settings" {
  name                 = "ZTNA settings"
  apply_to_entities    = [pfptmeta_group.ztna_group.id]
  max_devices_per_user = 5
  prohibited_os        = ["macOS", "iOS"]
}

resource "pfptmeta_user_settings" "login_settings" {
  name                = "Login settings"
  apply_on_org        = true
  sso_mandatory       = true
  mfa_required        = true
  allowed_factors     = ["SMS"]
  password_expiration = 30
}
