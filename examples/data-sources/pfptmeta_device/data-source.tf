data "pfptmeta_device" "my_device" {
  id = "dev-123"
}

output "my_device" {
  value = data.pfptmeta_device.my_device
}