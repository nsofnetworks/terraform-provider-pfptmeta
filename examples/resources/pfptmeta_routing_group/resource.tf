resource "pfptmeta_group" "group" {
  name = "group name"
}

resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service"
  mapped_service = "mapped.service.com"
}

resource "pfptmeta_routing_group" "routing_group" {
  name                = "routing group name"
  description         = "routing group description"
  mapped_elements_ids = [pfptmeta_network_element.mapped-service.id]
  sources             = [pfptmeta_group.group.id]
}