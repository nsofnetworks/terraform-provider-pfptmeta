resource "pfptmeta_file_scanning_rule" "file_scanning" {
  name             = "File scanning"
  description      = "File scanning rule"
  apply_to_org     = true
  block_file_types = ["exe"]
  malware          = "DOWNLOAD"
  priority         = 15
}