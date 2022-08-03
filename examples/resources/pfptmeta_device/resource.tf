resource "pfptmeta_device" "device" {
  name        = "my device"
  description = "some details about the device"
  owner_id    = "usr-abc123"
  tags = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
}