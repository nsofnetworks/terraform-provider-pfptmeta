resource "pfptmeta_tunnel" "tunnel1" {
  name        = "tunnel name1"
  description = "some details about the tunnel"
  gre_config {
    source_ups = ["198.51.100.15"]
  }
}
