resource "pfptmeta_pac_file" "managed_pac" {
  name         = "pac file"
  description  = "pac file description"
  apply_to_org = true
  priority     = 15
  type         = "managed"
  managed_content = {
    domains     = ["apple.co"]
    cloud_apps  = ["ca-zcR3wRh1mQLyD", "ca-YbU9LeY4bCAuV"]
    ip_networks = ["ipn-wRvGyyK1mPEop"]
  }
}

resource "pfptmeta_pac_file" "byo_pac" {
  name         = "pac file"
  description  = "pac file description"
  apply_to_org = true
  priority     = 15
  type         = "bring_your_own"
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

resource "pfptmeta_pac_file" "byo_pac_with_file_path" {
  name         = "pac file"
  description  = "pac file description"
  apply_to_org = true
  priority     = 15
  type         = "bring_your_own"
  content      = file("path/to/file")
}
