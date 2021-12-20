data "pfptmeta_alert" "alert" {
  id = "alr-123abc"
}

output "alert" {
  value = data.pfptmeta_alert.alert
}