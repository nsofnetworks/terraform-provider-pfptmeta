data "pfptmeta_certificate" "cert" {
  id = "crt-123abc"
}

output "certificate" {
  value = data.pfptmeta_certificate.cert
}