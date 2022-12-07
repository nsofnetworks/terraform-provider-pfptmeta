data "pfptmeta_cloud_app" "ca" {
  id = "ca-123abc"
}

output "content_category" {
  value = data.pfptmeta_cloud_app.ca
}