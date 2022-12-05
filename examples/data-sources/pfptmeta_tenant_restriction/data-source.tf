data "pfptmeta_tenant_restriction" "tr" {
  id = "cc-123abc"
}

output "tenant_restriction" {
  value = data.pfptmeta_tenant_restriction.tr
}