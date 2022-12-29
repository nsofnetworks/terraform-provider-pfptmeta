data "pfptmeta_ssl_bypass_rule" "rule" {
  id = "sbr-123abc"
}

output "rule" {
  value = data.pfptmeta_ssl_bypass_rule.rule
}