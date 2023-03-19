resource "pfptmeta_content_category" "web_mail_category" {
  name             = "Webmail Category"
  confidence_level = "LOW"
  types            = ["Web-based Email"]
  urls             = ["live.com", "outlook.com"]
}

resource "pfptmeta_url_filtering_rule" "isolate_web_mails" {
  name                         = "Isolate Web Mails"
  apply_to_org                 = true
  action                       = "ISOLATION"
  advanced_threat_protection   = false
  forbidden_content_categories = [pfptmeta_content_category.web_mail_category.id]
  priority                     = 92
  warn_ttl                     = 15
}
