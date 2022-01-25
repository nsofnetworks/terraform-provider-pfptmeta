resource "pfptmeta_notification_channel" "mail" {
  name = "mail-channel"
  email_config {
    recipients = ["user1@example.com", "user2@example.com"]
  }
}

resource "pfptmeta_log_streaming_access_bridge" "casb_log_stream" {
  name                  = "CASB log stream"
  description           = "log stream description"
  notification_channels = [pfptmeta_notification_channel.mail.id]
  export_logs           = ["traffic", "webfilter"]
  proofpoint_casb_config {
    region    = "EU"
    tenant_id = "tenant_70e4b2a567a24159ad6b495ba56b3620"
  }
}

resource "pfptmeta_log_streaming_access_bridge" "qradar_http_log_stream" {
  name                  = "QRadar HTTP log stream"
  description           = "log stream description"
  notification_channels = [pfptmeta_notification_channel.mail.id]
  export_logs           = ["api", "traffic", "security", "metaproxy", "webfilter"]
  qradar_http_config {
    certificate = file("/path/to/cert.pem")
    url         = "http://qradar.url:1234"
  }
}

resource "pfptmeta_log_streaming_access_bridge" "s3_log_stream" {
  name                  = "S3 log stream"
  description           = "log stream description"
  notification_channels = [pfptmeta_notification_channel.mail.id]
  export_logs           = ["api", "traffic", "security", "metaproxy", "webfilter"]
  s3_config {
    bucket   = "mybucket"
    compress = true
    prefix   = "mybucketstream"
  }
}

resource "pfptmeta_log_streaming_access_bridge" "splunk_http_log_stream" {
  name                  = "Splunk HTTP log stream"
  description           = "log stream description"
  notification_channels = [pfptmeta_notification_channel.mail.id]
  export_logs           = ["api", "traffic", "security", "metaproxy", "webfilter"]
  splunk_http_config {
    certificate         = <<-EOT
                          -----BEGIN CERTIFICATE-----
                          MIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIICCgKCAgEAlRuRnThUjU8/prwYxbty
                          WPT9pURI3lbsKMiB6Fn/VHOKE13p4D8xgOCADpdRagdT6n4etr9atzDKUSvpMtR3
                          CP5noNc97WiNCggBjVWhs7szEe8ugyqF23XwpHQ6uV1LKH50m92MbOWfCtjU9p/x
                          qhNpQQ1AZhqNy5Gevap5k8XzRmjSldNAFZMY7Yv3Gi+nyCwGwpVtBUwhuLzgNFK/
                          yDtw2WcWmUU7NuC8Q6MWvPebxVtCfVp/iQU6q60yyt6aGOBkhAX0LpKAEhKidixY
                          nP9PNVBvxgu3XZ4P36gZV6+ummKdBVnc3NqwBLu5+CcdRdusmHPHd5pHf4/38Z3/
                          6qU2a/fPvWzceVTEgZ47QjFMTCTmCwNt29cvi7zZeQzjtwQgn4ipN9NibRH/Ax/q
                          TbIzHfrJ1xa2RteWSdFjwtxi9C20HUkjXSeI4YlzQMH0fPX6KCE7aVePTOnB69I/
                          a9/q96DiXZajwlpq3wFctrs1oXqBp5DVrCIj8hU2wNgB7LtQ1mCtsYz//heai0K9
                          PhE4X6hiE0YmeAZjR0uHl8M/5aW9xCoJ72+12kKpWAa0SFRWLy6FejNYCYpkupVJ
                          yecLk/4L1W0l6jQQZnWErXZYe0PNFcmwGXy1Rep83kfBRNKRy5tvocalLlwXLdUk
                          AIU+2GKjyT3iMuzZxxFxPFMCAwEAAQ==
                          -----END CERTIFICATE-----
                          EOT
    publicly_accessible = false
    url                 = "https://splunk.url"
    token               = "token"
  }
}

resource "pfptmeta_log_streaming_access_bridge" "syslog_log_stream" {
  name                  = "Syslog log stream"
  description           = "log stream description"
  notification_channels = [pfptmeta_notification_channel.mail.id]
  export_logs           = ["api", "traffic", "security", "metaproxy", "webfilter"]
  syslog_config {
    host  = "syslog.hostname.com"
    port  = 518
    proto = "tcp"
  }
}