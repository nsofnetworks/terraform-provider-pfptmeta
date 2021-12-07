resource "pfptmeta_group" "new_group" {
  name        = "group name"
  description = "group description"
  expression  = "tag_name:tag_value OR platform:macOS"
}