data "pfptmeta_metaport_failover" "metaport_failover" {
  id = "mpc-123"
}

output "mapped_subnet" {
  value = data.pfptmeta_metaport_failover.metaport_failover
}