data "pfptmeta_protocol_group" "HTTP" {
  id = "pg-NKMzUnJzalxWZKe"
}

data "pfptmeta_protocol_group" "HTTPS" {
  name = "HTTPS"
}

output "protocol_group_by_id" {
  value = data.pfptmeta_protocol_group.HTTP
}

output "protocol_group_by_name" {
  value = data.pfptmeta_protocol_group.HTTPS
}