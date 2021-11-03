data "pfptmeta_mapped_subnet" "my_subnet" {
  id = "ne-123"
}

output "mapped_subnet" {
  value = data.pfptmeta_mapped_subnet.my_subnet
}