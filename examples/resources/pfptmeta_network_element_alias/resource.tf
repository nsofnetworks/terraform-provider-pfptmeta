resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service name"
  description    = "some details about the mapped service"
  mapped_service = "mapped.service.com"
}

resource "pfptmeta_network_element_alias" "alias" {
  network_element_id = pfptmeta_network_element.mapped-service.id
  alias              = "alias.test.com"
}