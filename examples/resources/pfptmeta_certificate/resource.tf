resource "pfptmeta_certificate" "cert" {
  name        = "certificate name"
  description = "certificate description"
  sans        = ["test.example.com"]
}