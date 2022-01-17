resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "mapped subnet name"
  mapped_subnets = ["0.0.0.0/0"]
}

resource "pfptmeta_metaport" "metaport" {
  name = "metaport name"
}

resource "pfptmeta_metaport_mapped_elements_attachment" "attachment" {
  metaport_id     = pfptmeta_metaport.metaport.id
  mapped_elements = [pfptmeta_network_element.mapped-subnet.id]
}