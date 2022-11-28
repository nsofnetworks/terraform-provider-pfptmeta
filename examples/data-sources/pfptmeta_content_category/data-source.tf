data "pfptmeta_content_category" "cc" {
  id = "cc-123abc"
}

output "content_category" {
  value = data.pfptmeta_content_category.cc
}