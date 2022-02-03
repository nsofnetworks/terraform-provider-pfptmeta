data "pfptmeta_metaport" "metaport_by_id" {
  id = "mp-123"
}

data "pfptmeta_metaport" "metaport_by_name" {
  name = "metaport name"
}

output "metaport_by_id" {
  value = data.pfptmeta_metaport.metaport_by_id
}

output "metaport_by_name" {
  value = data.pfptmeta_metaport.metaport_by_name
}