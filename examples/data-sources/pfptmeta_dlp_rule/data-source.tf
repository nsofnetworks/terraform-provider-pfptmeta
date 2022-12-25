data "pfptmeta_dlp_rule" "dlp" {
  id = "dlp-123abc"
}

output "dlp_rule" {
  value = data.pfptmeta_dlp_rule.dlp
}