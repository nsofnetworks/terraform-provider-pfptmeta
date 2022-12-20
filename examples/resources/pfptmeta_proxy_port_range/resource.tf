resource "pfptmeta_proxy_port_range" "proxy_port_range" {
  name        = "my port range"
  description = "some details about destination port ranges"
  proto       = "HTTP"
  from_port   = 20000
  to_port     = 20100
}
