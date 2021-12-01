data "pfptmeta_mapped_host" "mapped-host" {
  network_element_id = "ne-123"
  name               = "host.test.com"
}

output "alias" {
  value = data.pfptmeta_mapped_host.mapped-host
}