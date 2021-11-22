data "pfptmeta_metaport" "metaport1" {
  id = "mp-123"
}

output "mapped_subnet" {
  value = data.pfptmeta_metaport.metaport1
}