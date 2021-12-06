data "pfptmeta_metaport_cluster" "metaport_cluster" {
  id = "mpc-123"
}

output "mapped_subnet" {
  value = data.pfptmeta_metaport_cluster.metaport_cluster
}