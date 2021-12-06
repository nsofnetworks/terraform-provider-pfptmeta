data "pfptmeta_mapped_domain" "mapped-domain" {
  network_element_id = "ne-123"
  name               = "test.com"
}

output "mapped_domain" {
  value = data.pfptmeta_mapped_domain.mapped-domain
}