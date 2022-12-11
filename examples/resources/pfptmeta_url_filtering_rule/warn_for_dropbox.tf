data "pfptmeta_catalog_app" "dropbox" {
  name     = "Dropbox"
  category = "Collaboration"
}

resource "pfptmeta_cloud_app" "dropbox_personal" {
  name        = "Dropbox Personal"
  app         = data.pfptmeta_catalog_app.dropbox.id
  tenant_type = "Personal"
}

resource "pfptmeta_url_filtering_rule" "warn_for_dropbox" {
  name         = "Warn For Dropbox"
  apply_to_org = true
  action       = "WARN"
  cloud_apps   = [pfptmeta_cloud_app.dropbox_personal.id]
  priority     = 88
  warn_ttl     = 15
}