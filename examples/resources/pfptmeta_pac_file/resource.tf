resource "pfptmeta_pac_file" "pac" {
  name         = "pac file"
  description  = "pac file description"
  apply_to_org = true
  priority     = 15
  content      = <<EOF
function FindProxyForURL(url, host) {
// Don't proxy specific hostnames
if (
  shExpMatch(host, "*.apple.co")
  )
  return "DIRECT";
return "PROXY 127.0.0.1:43443";
}
EOF
}

resource "pfptmeta_pac_file" "pac_with_file_path" {
  name         = "pac file"
  description  = "pac file description"
  apply_to_org = true
  priority     = 15
  content      = file("path/to/file")
}