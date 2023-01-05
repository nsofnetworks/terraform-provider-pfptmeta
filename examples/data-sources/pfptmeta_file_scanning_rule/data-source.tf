data "pfptmeta_file_scanning_rule" "fsr" {
  id = "fsr-123abc"
}

output "fsr" {
  value = data.pfptmeta_file_scanning_rule.fsr
}