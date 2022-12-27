data "pfptmeta_pac_file" "pac_file" {
  id = "alr-123abc"
}

output "pac_file" {
  value = data.pfptmeta_pac_file.pac_file
}