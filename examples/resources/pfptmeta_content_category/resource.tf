resource "pfptmeta_content_category" "cc" {
  name                      = "content category"
  description               = "content category description"
  confidence_level          = "HIGH"
  forbid_uncategorized_urls = true
  types                     = ["News and Media", "Sports"]
  urls                      = ["192.6.6.5", "ynet.co.il"]
}

