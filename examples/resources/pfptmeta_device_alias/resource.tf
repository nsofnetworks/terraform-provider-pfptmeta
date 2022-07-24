resource "pfptmeta_device" "device" {
  name        = "a name for the device"
  description = "some details about the device"
  owner_id    = "usr-abc123"
}

resource "pfptmeta_device_alias" "alias" {
  device_id = pfptmeta_device.device.id
  alias     = "alias.test.com"
}