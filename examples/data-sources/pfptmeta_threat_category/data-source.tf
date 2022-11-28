data "pfptmeta_threat_category" "tc" {
  id = "tc-123abc"
}

output "threat_category" {
  value = data.pfptmeta_threat_category.tc
}