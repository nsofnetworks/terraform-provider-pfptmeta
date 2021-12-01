data "pfptmeta_network_element_alias" "alias" {
  network_element_id = "ne-123"
  alias              = "test.com"
}

output "alias" {
  value = data.pfptmeta_network_element_alias.alias
}