resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "mapped service name"
  description    = "some details about the mapped service"
  mapped_subnets = ["0.0.0.0/0", "10.20.30.0/24"]
}

resource "pfptmeta_mapped_host" "mapped-host" {
  network_element_id = pfptmeta_network_element.mapped-subnet.id
  mapped_host        = "host.test.com"
  name               = "host.test.com"
}