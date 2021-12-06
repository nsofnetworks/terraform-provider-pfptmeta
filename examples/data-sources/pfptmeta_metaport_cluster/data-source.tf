data "pfptmeta_metaport_cluster" "metaport_cluster" {
  id = "mpc-123"
}

output "metaport_cluster" {
  value = data.pfptmeta_metaport_cluster.metaport_cluster
}