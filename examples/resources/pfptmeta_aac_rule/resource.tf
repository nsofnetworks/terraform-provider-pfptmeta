resource "pfptmeta_aac_rule" "aac_rule" {
  name             = "aac rule name"
  description      = "aac rule description"
  enabled          = true
  priority         = 555
  action           = "allow"
  app_ids          = ["app-abcd1234"]
  sources          = ["usr-abcd1234"]
  certificate_id   = "crt-abcd1234"
  suspicious_login = "safe"
  locations        = ["US", "IL"]
  ip_reputations   = ["tor", "vpn"]
}