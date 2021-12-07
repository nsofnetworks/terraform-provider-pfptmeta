data "pfptmeta_group" "group_by_id" {
  id = "grp-123abc"
}

data "pfptmeta_group" "group_by_name" {
  name = "dev"
}

output "group_by_id" {
  value = data.pfptmeta_group.group_by_id
}

output "group_by_name" {
  value = data.pfptmeta_group.group_by_name
}