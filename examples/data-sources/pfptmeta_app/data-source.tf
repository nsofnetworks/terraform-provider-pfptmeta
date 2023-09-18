data "pfptmeta_app" "app" {
  id       = "app-1234abcd"
  protocol = "SAML"
}

output "app" {
  value = data.pfptmeta_app.app
}