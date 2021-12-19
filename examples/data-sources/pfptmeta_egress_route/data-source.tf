data "pfptmeta_egress_route" "egress_route" {
  id = "er-123abc"
}

output "egress_route" {
  value = data.pfptmeta_egress_route.egress_route
}