data "pfptmeta_catalog_app" "g_suite" {
  name     = "G Suite"
  category = "Utilities"
}

resource "pfptmeta_cloud_app" "g_suite" {
  name = "salesforce"
  app  = data.pfptmeta_catalog_app.g_suite.id
  urls = [
    "gmail.com",
    "googleusercontent.com",
    "clients6.google.com",
  ]
}

data "pfptmeta_catalog_app" "dropbox" {
  name     = "Dropbox"
  category = "Collaboration"
}

resource "pfptmeta_cloud_app" "dropbox_personal" {
  name        = "Dropbox Personal"
  app         = data.pfptmeta_catalog_app.dropbox.id
  tenant_type = "Personal"
}

resource "pfptmeta_dlp_rule" "default_rule" {
  name                     = "default rule"
  description              = "Block Dropbox and Google Drive"
  apply_to_org             = true
  action                   = "BLOCK"
  alert_level              = "MEDIUM"
  all_supported_file_types = true
  file_parts               = ["CONTENT"]
  cloud_apps               = [pfptmeta_cloud_app.g_suite.id, pfptmeta_cloud_app.dropbox_personal.id]
  user_actions             = ["UPLOAD", "DOWNLOAD"]
  priority                 = 15
}
