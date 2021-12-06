data "pfptmeta_metaport_failover" "metaport_failover" {
  id = "mpc-123"
}

output "metaport_failover" {
  value = data.pfptmeta_metaport_failover.metaport_failover
}