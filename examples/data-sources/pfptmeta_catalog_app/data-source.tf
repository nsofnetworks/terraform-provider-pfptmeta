data "pfptmeta_catalog_app" "catalog_app" {
  name     = "Google"
  category = "Property Management"
}

output "catalog_app" {
  value = data.pfptmeta_catalog_app.catalog_app
}