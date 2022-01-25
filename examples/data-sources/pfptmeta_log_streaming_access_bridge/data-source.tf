data "pfptmeta_log_streaming_access_bridge" "log_streaming" {
  id = "crt-123abc"
}

output "log_streaming" {
  value = data.pfptmeta_log_streaming_access_bridge.log_streaming
}