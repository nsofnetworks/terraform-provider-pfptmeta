---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Resource pfptmeta_log_streaming_access_bridge - terraform-provider-pfptmeta"
subcategory: "Administration"
description: |-
  Proofpoint log streaming allows you to send log data from the Proofpoint NaaS tenant to a third-party service in real time.
  Then, the log data can be stored and analyzed using a Security Information and Event Management (SIEM) solution.
  Proofpoint supports the following log streaming standards:
  Splunk-compatible HTTP Event Collector (HEC), enabling you to send data over HTTP (or HTTPS) directly to Splunk Enterprise or Splunk Cloud from Proofpoint NaaS. Splunk HEC is token-based, eliminating the need to hard-code your Splunk credentials into Proofpoint-provided applications.IBM QRadar-compatible HTTP for collecting flow and event data from all of the log sources that are supported in your on-premises or cloud deployment.Amazon S3 service for direct streaming of the customer’s tenant logs to an AWS S3 bucket.Syslog Common Event Format (CEF), an open-source log management standard. CEF allows third parties to create their own device schemas that are compatible with industry-standard methods for normalizing security events.Proofpoint CASB - Proofpoint CASB service that accepts traffic or web security logs for subsequent shadow IT processing.
  This integration enables organizations to govern user access to both IT-authorized and unauthorized apps (also known as shadow IT)
  You can use pfptmetanotificationchannel https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/notification_channel for receiving alerts on events from log streaming service.
  When the log streamer status changes to one of the non-running states (error, suspended, stopped) a notification event is triggered.
---

# Resource (pfptmeta_log_streaming_access_bridge)

Proofpoint log streaming allows you to send log data from the Proofpoint NaaS tenant to a third-party service in real time. 
Then, the log data can be stored and analyzed using a Security Information and Event Management (SIEM) solution.
Proofpoint supports the following log streaming standards:

- Splunk-compatible HTTP Event Collector (HEC), enabling you to send data over HTTP (or HTTPS) directly to Splunk Enterprise or Splunk Cloud from Proofpoint NaaS. Splunk HEC is token-based, eliminating the need to hard-code your Splunk credentials into Proofpoint-provided applications.
- IBM QRadar-compatible HTTP for collecting flow and event data from all of the log sources that are supported in your on-premises or cloud deployment.
- Amazon S3 service for direct streaming of the customer’s tenant logs to an AWS S3 bucket.
- Syslog Common Event Format (CEF), an open-source log management standard. CEF allows third parties to create their own device schemas that are compatible with industry-standard methods for normalizing security events.
- Proofpoint CASB - Proofpoint CASB service that accepts traffic or web security logs for subsequent shadow IT processing.
This integration enables organizations to govern user access to both IT-authorized and unauthorized apps (also known as shadow IT)

You can use [**pfptmeta_notification_channel**](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/notification_channel) for receiving alerts on events from log streaming service.
When the log streamer status changes to one of the non-running states (error, suspended, stopped) a notification event is triggered.

## Example Usage

```terraform
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `export_logs` (List of String) Enum: `api` `traffic` `security` `metaproxy` `webfilter`.
- `name` (String)

### Optional

- `description` (String)
- `enabled` (Boolean)
- `notification_channels` (List of String) Notification channel IDs to which an alert will be sent if the log streaming service becomes unavailable or the endpoint is unreachable.
- `proofpoint_casb_config` (Block List, Max: 1) Configuration for log streaming to Proofpoint CASB for shadow IT processing. (see [below for nested schema](#nestedblock--proofpoint_casb_config))
- `qradar_http_config` (Block List, Max: 1) Configuration for log streaming to IBM QRadar platform. (see [below for nested schema](#nestedblock--qradar_http_config))
- `s3_config` (Block List, Max: 1) Configuration for log streaming to an Amazon S3 bucket. (see [below for nested schema](#nestedblock--s3_config))
- `splunk_http_config` (Block List, Max: 1) Configuration for log streaming to Self-Hosted / cloud Splunk. see [here](https://help.metanetworks.com/knowledgebase/log_streaming_for_splunk_self_hosted/#configuring-splunk-http-event-collector) for instructions on how to enable HTTP Event Collector on Self-Hosted Instance, and [here](https://help.metanetworks.com/knowledgebase/log_streaming_for_splunk_cloud/#configuring-splunk-http-event-collector) for instructions on how to enable HTTP Event Collector on Cloud Instance. (see [below for nested schema](#nestedblock--splunk_http_config))
- `syslog_config` (Block List, Max: 1) Configuration for log streaming in Syslog Common Event Format (CEF). (see [below for nested schema](#nestedblock--syslog_config))

### Read-Only

- `id` (String) The ID of this resource.
- `status` (String)
- `status_description` (String)

<a id="nestedblock--proofpoint_casb_config"></a>
### Nested Schema for `proofpoint_casb_config`

Required:

- `region` (String) Tenant region in Proofpoint CASB system. ENUM: `EU`, `US`.
- `tenant_id` (String) Your organization CASB tenant ID in the tenant_ID format. Valid tenant ID format is *tenant<32 hex characters>*.


<a id="nestedblock--qradar_http_config"></a>
### Nested Schema for `qradar_http_config`

Required:

- `url` (String) QRadar server URL.

Optional:

- `certificate` (String, Sensitive) Base64 root CA certificate.


<a id="nestedblock--s3_config"></a>
### Nested Schema for `s3_config`

Required:

- `bucket` (String) Name of your AWS S3 bucket. The bucket must have proper writing permissions and access rights for the Proofpoint logging.
- `compress` (Boolean) Defines whether to compress log objects using gzip or not.

Optional:

- `prefix` (String) Shared name prefix for destination object in your AWS S3 bucket.


<a id="nestedblock--splunk_http_config"></a>
### Nested Schema for `splunk_http_config`

Required:

- `token` (String, Sensitive) Token code generated by Splunk HEC.
- `url` (String) Splunk server HTTP event collector URI. This field accepts an FQDN value only.

Optional:

- `certificate` (String, Sensitive) Base 64 server for the root CA certificate.
- `publicly_accessible` (Boolean) Whether the Splunk instance URL endpoint is publicly available.


<a id="nestedblock--syslog_config"></a>
### Nested Schema for `syslog_config`

Required:

- `host` (String) SIEM destination FQDN.
- `port` (Number) TCP port for log data input.
- `proto` (String) ENUM: `tcp`, `udp`.
