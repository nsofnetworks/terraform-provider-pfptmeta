resource "pfptmeta_certificate" "managed_cert" {
  name        = "certificate name"
  description = "certificate description"
  sans        = ["test.example.com"]
}

resource "pfptmeta_certificate" "byo_cert" {
  name        = "certificate name"
  description = "certificate description"
  certificate = file("my_cert.pem")
}
