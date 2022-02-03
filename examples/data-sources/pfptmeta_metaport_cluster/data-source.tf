data "pfptmeta_metaport_cluster" "metaport_cluster_by_id" {
  id = "mpc-123"
}

output "metaport_cluster_by_id" {
  value = data.pfptmeta_metaport_cluster.metaport_cluster_by_id
}

data "pfptmeta_metaport_cluster" "metaport_cluster_by_name" {
  name = "metaport cluster name"
}

output "metaport_cluster_by_name" {
  value = data.pfptmeta_metaport_cluster.metaport_cluster_by_name
}