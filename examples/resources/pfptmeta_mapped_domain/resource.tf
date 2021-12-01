resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "mapped service name"
  description    = "some details about the mapped service"
  mapped_subnets = ["0.0.0.0/0", "10.20.30.0/24"]
}

resource "pfptmeta_mapped_domain" "mapped-domain" {
  network_element_id = pfptmeta_network_element.mapped-subnet.id
  mapped_domain      = "test.com"
  name               = "test1.com"
}