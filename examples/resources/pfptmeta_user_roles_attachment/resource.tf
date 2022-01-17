resource "pfptmeta_user" "user" {
  given_name  = "John"
  family_name = "Smith"
  email       = "john.smith@example.com"
}

resource "pfptmeta_role" "metaport_role" {
  name       = "metaport role"
  privileges = ["metaports:read", "metaports:write"]
}

resource "pfptmeta_role" "network_element_role" {
  name       = "network element role"
  privileges = ["network_elements:read", "network_elements:write"]
}

resource "pfptmeta_user_roles_attachment" "attachment" {
  user_id = pfptmeta_user.user.id
  roles   = [pfptmeta_role.metaport_role.id, pfptmeta_role.network_element_role.id]
}