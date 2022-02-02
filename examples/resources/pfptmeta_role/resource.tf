locals {
  org_id = "org123abc"
}

resource "pfptmeta_role" "read_only" {
  name                = "all read only"
  description         = "role with all read privileges"
  apply_to_orgs       = [local.org_id]
  all_read_privileges = true
}

resource "pfptmeta_role" "admin_role" {
  name                 = "all read only"
  description          = "role with all read privileges"
  all_read_privileges  = true
  all_write_privileges = true
  all_suborgs          = true
}

resource "pfptmeta_role" "with_privileges" {
  name          = "with privs"
  apply_to_orgs = [local.org_id]
  privileges    = ["metaports:read", "metaports:write"]
}