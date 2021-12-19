data "pfptmeta_group" "group" {
  name = "group-name"
}

resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "mapped subnet name"
  mapped_subnets = ["10.20.30.0/24"]
}

data "pfptmeta_location" "new_york" {
  name = "LGA"
}

resource "pfptmeta_egress_route" "via_mapped_subnet" {
  name         = "egress route name"
  description  = "egress route description"
  destinations = ["example.com"]
  sources      = [data.pfptmeta_group.group.id]
  via          = pfptmeta_network_element.mapped-subnet.id
}

resource "pfptmeta_egress_route" "via_region" {
  name         = "example"
  destinations = ["example.com"]
  sources      = [data.pfptmeta_group.group.id]
  via          = data.pfptmeta_location.new_york.name
}

resource "pfptmeta_egress_route" "DIRECT" {
  name         = "example"
  destinations = ["example.com"]
  sources      = [data.pfptmeta_group.group.id]
  via          = "DIRECT"
}