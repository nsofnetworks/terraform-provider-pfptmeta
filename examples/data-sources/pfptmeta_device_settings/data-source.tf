data "pfptmeta_device_settings" "settings" {
  id = "ds-123abc"
}

output "settings" {
  value = data.pfptmeta_device_settings.settings
}