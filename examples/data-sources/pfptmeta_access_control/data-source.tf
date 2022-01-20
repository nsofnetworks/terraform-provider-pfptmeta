data "pfptmeta_access_control" "access" {
  id = "ac-123abc"
}

output "access" {
  value = data.pfptmeta_access_control.access
}