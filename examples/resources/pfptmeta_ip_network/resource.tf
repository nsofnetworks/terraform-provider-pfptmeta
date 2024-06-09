resource "pfptmeta_ip_network" "in" {
  name        = "ip network name"
  description = "ip network description"
  cidrs       = ["0.0.0.0/32"]
  countries   = ["US"]
}
