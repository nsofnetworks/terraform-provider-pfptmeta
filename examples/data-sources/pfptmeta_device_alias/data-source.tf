data "pfptmeta_device_alias" "alias" {
  device_id = "dev-123"
  alias     = "test.com"
}

output "alias" {
  value = data.pfptmeta_device_alias.alias
}