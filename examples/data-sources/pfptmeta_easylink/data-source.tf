data "pfptmeta_easylink" "easylink" {
  id = "el-123abc"
}

output "easylink" {
  value = data.pfptmeta_easylink.easylink
}