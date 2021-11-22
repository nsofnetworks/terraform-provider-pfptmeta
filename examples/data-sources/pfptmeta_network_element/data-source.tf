data "pfptmeta_network_element" "my_subnet" {
  id = "ne-123"
}

output "mapped_subnet" {
  value = data.pfptmeta_network_element.my_subnet
}