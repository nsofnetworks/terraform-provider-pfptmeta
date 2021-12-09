resource "pfptmeta_group" "group" {
  name = "admins"
}

resource "pfptmeta_role" "metaport_role" {
  name       = "metaport role"
  privileges = ["metaports:read", "metaports:write"]
}

resource "pfptmeta_role" "network_element_role" {
  name       = "network element role"
  privileges = ["network_elements:read", "network_elements:write"]
}

resource "pfptmeta_group_roles_attachment" "attachment" {
  group_id = pfptmeta_group.group.id
  roles    = [pfptmeta_role.metaport_role.id, pfptmeta_role.network_element_role.id]
}