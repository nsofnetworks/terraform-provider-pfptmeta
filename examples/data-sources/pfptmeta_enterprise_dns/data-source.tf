data "pfptmeta_enterprise_dns" "enterprise_dns" {
  id = "ed-123"
}

output "enterprise_dns" {
  value = data.pfptmeta_enterprise_dns.enterprise_dns
}