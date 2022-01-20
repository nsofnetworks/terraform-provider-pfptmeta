resource "pfptmeta_access_control" "access" {
  name           = "access control name"
  description    = "access control description"
  apply_to_org   = true
  allowed_routes = ["192.0.2.1/24", "192.168.0.1/24"]
}