data "pfptmeta_metaport" "metaport1" {
  id = "mp-123"
}

output "metaport" {
  value = data.pfptmeta_metaport.metaport1
}