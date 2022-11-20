resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "ms"
  mapped_subnets = ["0.0.0.0/0"]
}

resource "pfptmeta_metaport" "metaport1" {
  name        = "metaport name"
  description = "metaport description"
}

resource "pfptmeta_metaport" "metaport2" {
  name        = "metaport name"
  description = "metaport description"
}

resource "pfptmeta_metaport_cluster" "metaport_cluster" {
  name        = "metaport cluster name"
  description = "metaport cluster description"
  metaports   = [pfptmeta_metaport.metaport1.id, pfptmeta_metaport.metaport2.id]
}

resource "pfptmeta_metaport_cluster_mapped_elements_attachment" "attachment" {
  metaport_cluster_id = pfptmeta_metaport_cluster.metaport_cluster.id
  mapped_elements     = [pfptmeta_network_element.mapped-subnet.id]
}