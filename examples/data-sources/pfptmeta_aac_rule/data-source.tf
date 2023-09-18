data "pfptmeta_aac_rule" "aac_rule" {
  id = "arl-XWU8BQmGbevn7"
}

output "aac_rule" {
  value = data.pfptmeta_aac_rule.aac_rule
}