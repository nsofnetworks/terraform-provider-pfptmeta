data "pfptmeta_notification_channel" "channel" {
  id = "nc-123abc"
}

output "mapped_subnet" {
  value = data.pfptmeta_notification_channel.channel
}