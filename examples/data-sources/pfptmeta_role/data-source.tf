data "pfptmeta_role" "admin_by_id" {
  id = "rol-Xqrzun95v8RA59E"
}

data "pfptmeta_role" "admin_by_name" {
  name = "admin"
}

output "role_by_id" {
  value = data.pfptmeta_role.admin_by_id
}

output "role_by_name" {
  value = data.pfptmeta_role.admin_by_name
}