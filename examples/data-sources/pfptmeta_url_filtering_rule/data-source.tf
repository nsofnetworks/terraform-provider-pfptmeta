data "pfptmeta_url_filtering_rule" "ufr" {
  id = "ufr-123abc"
}

output "catalog_app" {
  value = data.pfptmeta_url_filtering_rule.ufr
}