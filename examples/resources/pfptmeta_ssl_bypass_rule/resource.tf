resource "pfptmeta_ssl_bypass_rule" "rule" {
  name                      = "rule name"
  description               = "rule description"
  apply_to_org              = true
  bypass_uncategorized_urls = false
  content_types             = ["Abortion"]
  domains                   = [".youtube.com"]
  priority                  = 15
}
