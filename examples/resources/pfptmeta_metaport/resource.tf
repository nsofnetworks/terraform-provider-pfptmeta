resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "mapped subnet name"
  mapped_subnets = ["0.0.0.0/0"]
}

resource "pfptmeta_metaport" "metaport1" {
  name            = "metaport name1"
  description     = "some details about the metaport"
  mapped_elements = [pfptmeta_network_element.mapped-subnet.id]
  allow_support   = false
}