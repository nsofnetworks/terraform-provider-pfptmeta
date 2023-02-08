data "pfptmeta_tunnel" "tunnel_by_id" {
  id = "tun-123"
}

data "pfptmeta_tunnel" "tunnel_by_name" {
  name = "tunnel name"
}

output "tunnel_by_id" {
  value = data.pfptmeta_tunnel.tunnel_by_id
}

output "tunnel_by_name" {
  value = data.pfptmeta_tunnel.tunnel_by_name
}
