data "pfptmeta_idp" "idp" {
  id = "idp-1234abcd"
}

output "idp" {
  value = data.pfptmeta_idp.idp
}