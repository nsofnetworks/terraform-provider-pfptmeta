data "pfptmeta_user_settings" "settings" {
  id = "ds-123abc"
}

output "settings" {
  value = data.pfptmeta_user_settings.settings
}