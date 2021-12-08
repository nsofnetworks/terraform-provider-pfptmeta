data "pfptmeta_routing_group" "routing_group" {
  id = "rg-abc123"
}

output "routing_group" {
  value = data.pfptmeta_routing_group.routing_group
}