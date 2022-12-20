data "pfptmeta_ppr" "my_proxy_port_range" {
  id = "ppr-1a2b3c"
}

output "my_proxy_port_range" {
  value = data.pfptmeta_device.my_proxy_port_range
}