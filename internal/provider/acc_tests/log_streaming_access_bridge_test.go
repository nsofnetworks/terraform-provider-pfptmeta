package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	logStreamDependencies = `
resource "pfptmeta_notification_channel" "mail" {
  name = "mail-channel"
  email_config {
    recipients = ["user1@example.com", "user2@example.com"]
  }
}
`
	logStreamAccessBridgeCasbStep1 = `
resource "pfptmeta_log_streaming_access_bridge" "casb_log_stream" {
  name                  = "CASB log stream"
  description           = "log stream description"
  notification_channels = [pfptmeta_notification_channel.mail.id]
  siem_config {
    export_logs = ["traffic", "webfilter"]
    proofpoint_casb_config {
      region    = "EU"
      tenant_id = "tenant_70e4b2a567a24159ad6b495ba56b3620"
    }
  }
}
`
	logStreamAccessBridgeCasbStep2 = `
resource "pfptmeta_log_streaming_access_bridge" "casb_log_stream" {
  name                  = "CASB log stream1"
  description           = "log stream description1"
  notification_channels = [pfptmeta_notification_channel.mail.id]
  siem_config {
    export_logs = ["webfilter", "traffic"]
    proofpoint_casb_config {
      region    = "US"
      tenant_id = "tenant_70e4b2a567a24159ad6b495ba56b3621"
    }
  }
}
`
	logStreamAccessBridgeQRadar = `
resource "pfptmeta_log_streaming_access_bridge" "qradar_http_log_stream" {
  name                  = "QRadar HTTP log stream"
  description           = "log stream description"
  notification_channels = [pfptmeta_notification_channel.mail.id]
  siem_config {
    export_logs = ["api", "traffic", "security", "metaproxy", "webfilter"]
    qradar_http_config {
	  certificate  = <<-EOT
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
      url         = "http://qradar.url:1234"
    }
  }
}
`
	logStreamAccessBridgeS3 = `
resource "pfptmeta_log_streaming_access_bridge" "s3_log_stream" {
  name                  = "S3 log stream"
  description           = "log stream description"
  notification_channels = [pfptmeta_notification_channel.mail.id]
  siem_config {
    export_logs = ["api", "traffic", "security", "metaproxy", "webfilter"]
    s3_config {
      bucket   = "mybucket"
      compress = true
      prefix   = "mybucketstream"
    }
  }
}
`

	logStreamAccessBridgeSplunk = `
	resource "pfptmeta_log_streaming_access_bridge" "splunk_http_log_stream" {
  name                  = "Splunk HTTP log stream"
  description           = "log stream description"
  notification_channels = [pfptmeta_notification_channel.mail.id]
  siem_config {
    export_logs = ["api", "traffic", "security", "metaproxy", "webfilter"]
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
      token               = "token"
      url                 = "https://splunk.url"
    }
  }
}
`
	logStreamAccessBridgeSyslog = `
resource "pfptmeta_log_streaming_access_bridge" "syslog_log_stream" {
  name                  = "Syslog log stream"
  description           = "log stream description"
  notification_channels = [pfptmeta_notification_channel.mail.id]
  siem_config {
    export_logs = ["api", "traffic", "security", "metaproxy", "webfilter"]
    syslog_config {
      host  = "syslog.hostname.com"
      port  = 518
      proto = "tcp"
    }
  }
}`
	logStreamAccessBridgeDataSource = `
data "pfptmeta_log_streaming_access_bridge" "log_stream" {
  id = pfptmeta_log_streaming_access_bridge.casb_log_stream.id
}`
)

func TestAccResourceLogStreamAccessBridge(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("pfptmeta_log_streaming_access_bridge", "v1/access_bridges"),
		Steps: []resource.TestStep{
			{
				Config: logStreamDependencies +
					logStreamAccessBridgeCasbStep1 +
					logStreamAccessBridgeQRadar +
					logStreamAccessBridgeS3 +
					logStreamAccessBridgeSplunk +
					logStreamAccessBridgeSyslog,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream", "id", regexp.MustCompile("^ab-.+$")),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream", "name", "CASB log stream"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream", "description", "log stream description"),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream", "notification_channels.0",
						"pfptmeta_notification_channel.mail", "id"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream", "siem_config.0.export_logs.0", "traffic"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream", "siem_config.0.export_logs.1", "webfilter"),

					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream",
						"siem_config.0.proofpoint_casb_config.0.region", "EU"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream",
						"siem_config.0.proofpoint_casb_config.0.tenant_id", "tenant_70e4b2a567a24159ad6b495ba56b3620"),

					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.qradar_http_log_stream",
						"siem_config.0.qradar_http_config.0.url", "http://qradar.url:1234"),

					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.s3_log_stream",
						"siem_config.0.s3_config.0.bucket", "mybucket"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.s3_log_stream",
						"siem_config.0.s3_config.0.compress", "true"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.s3_log_stream",
						"siem_config.0.s3_config.0.prefix", "mybucketstream"),

					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.splunk_http_log_stream",
						"siem_config.0.splunk_http_config.0.publicly_accessible", "false"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.splunk_http_log_stream",
						"siem_config.0.splunk_http_config.0.token", "token"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.splunk_http_log_stream",
						"siem_config.0.splunk_http_config.0.url", "https://splunk.url"),

					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.syslog_log_stream",
						"siem_config.0.syslog_config.0.host", "syslog.hostname.com"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.syslog_log_stream",
						"siem_config.0.syslog_config.0.port", "518"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.syslog_log_stream",
						"siem_config.0.syslog_config.0.proto", "tcp"),
				),
			},
			{
				Config: logStreamDependencies + logStreamAccessBridgeCasbStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream", "id", regexp.MustCompile("^ab-.+$")),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream", "name", "CASB log stream1"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream", "description", "log stream description1"),
					resource.TestCheckResourceAttrPair(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream", "notification_channels.0",
						"pfptmeta_notification_channel.mail", "id"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream", "siem_config.0.export_logs.0", "webfilter"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream", "siem_config.0.export_logs.1", "traffic"),

					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream",
						"siem_config.0.proofpoint_casb_config.0.region", "US"),
					resource.TestCheckResourceAttr(
						"pfptmeta_log_streaming_access_bridge.casb_log_stream",
						"siem_config.0.proofpoint_casb_config.0.tenant_id", "tenant_70e4b2a567a24159ad6b495ba56b3621"),
				),
			},
		},
	})
}

func TestAccDataSourceLogStreamAccessBridge(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("pfptmeta_log_streaming_access_bridge", "v1/access_bridges"),
		Steps: []resource.TestStep{
			{
				Config: logStreamDependencies + logStreamAccessBridgeCasbStep1 + logStreamAccessBridgeDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_log_streaming_access_bridge.log_stream", "id", regexp.MustCompile("^ab-.+$")),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_log_streaming_access_bridge.log_stream", "name", "CASB log stream"),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_log_streaming_access_bridge.log_stream", "description", "log stream description"),
					resource.TestCheckResourceAttrPair(
						"data.pfptmeta_log_streaming_access_bridge.log_stream", "notification_channels.0",
						"pfptmeta_notification_channel.mail", "id"),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_log_streaming_access_bridge.log_stream", "siem_config.0.export_logs.0", "traffic"),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_log_streaming_access_bridge.log_stream", "siem_config.0.export_logs.1", "webfilter"),

					resource.TestCheckResourceAttr(
						"data.pfptmeta_log_streaming_access_bridge.log_stream",
						"siem_config.0.proofpoint_casb_config.0.region", "EU"),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_log_streaming_access_bridge.log_stream",
						"siem_config.0.proofpoint_casb_config.0.tenant_id", "tenant_70e4b2a567a24159ad6b495ba56b3620"),
				),
			},
		},
	})
}
