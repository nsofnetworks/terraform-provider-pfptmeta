data "pfptmeta_trusted_network" "network" {
  id = "tn-123abc"
}

output "network" {
  value = data.pfptmeta_trusted_network.network
}