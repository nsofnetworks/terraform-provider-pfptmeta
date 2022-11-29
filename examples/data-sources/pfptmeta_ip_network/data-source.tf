data "pfptmeta_ip_network" "in" {
  id = "ipn-123abc"
}

output "ip_network" {
  value = data.pfptmeta_ip_network.in
}