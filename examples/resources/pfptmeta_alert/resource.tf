resource "pfptmeta_notification_channel" "channel" {
  name = "mail-channel"
  email_config {
    recipients = ["user1@example.com", "user2@example.com"]
  }
}

resource "pfptmeta_alert" "spike_condition_alert" {
  name           = "spike"
  description    = "alert on metaport disconnection"
  channels       = [pfptmeta_notification_channel.channel.id]
  group_by       = "src_id"
  notify_message = "metaport disconnected"
  query_text     = "event:keepalive AND src_type:MetaPort"
  source_type    = "traffic_audit"
  spike_condition {
    spike_ratio = 100
    spike_type  = "down"
    time_diff   = 5
  }
  window = 5
}


resource "pfptmeta_alert" "threshold_condition_alert" {
  name           = "threshold"
  description    = "alert on new device onboard"
  channels       = [pfptmeta_notification_channel.channel.id]
  notify_message = "device onboarded"
  query_text     = "action:CREATE AND resource_type:Device"
  source_type    = "api_audit"
  threshold_condition {
    formula   = "count"
    op        = "greater"
    threshold = 0
  }
  window = 60
}