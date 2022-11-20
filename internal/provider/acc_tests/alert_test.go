package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	notificationChannelConf = `
resource "pfptmeta_notification_channel" "channel" {
  name = "mail-channel"
  email_config {
    recipients = ["user1@example.com", "user2@example.com"]
  }
}`
	spikeConditionAlert = `
resource "pfptmeta_alert" "alert" {
  name           = "alert-name"
  description    = "alert-description"
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
}`
	spikeConditionAlertStep2 = `
resource "pfptmeta_alert" "alert" {
  name           = "alert-name1"
  description    = "alert-description1"
  channels       = [pfptmeta_notification_channel.channel.id]
  group_by       = "src_id"
  notify_message = "metaport disconnected1"
  query_text     = "event:keepalive"
  source_type    = "traffic_audit"
  spike_condition {
    spike_ratio = 99
    spike_type  = "down"
    time_diff   = 60
  }
  window = 60
}`
	thresholdConditionAlert = `
resource "pfptmeta_alert" "alert2" {
  name           = "alert-name1"
  description    = "alert-description1"
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
}`
	alertDataSource = `
data "pfptmeta_alert" "alert" {
  id = pfptmeta_alert.alert.id
}
`
)

func TestAccResourceAlert(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("alert", "v1/alerts"),
		Steps: []resource.TestStep{
			{
				Config: notificationChannelConf + spikeConditionAlert,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_alert.alert", "id", regexp.MustCompile("^alr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "name", "alert-name"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "description", "alert-description"),
					resource.TestMatchResourceAttr("pfptmeta_alert.alert", "channels.0", regexp.MustCompile("^nch-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "group_by", "src_id"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "notify_message", "metaport disconnected"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "query_text", "event:keepalive AND src_type:MetaPort"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "source_type", "traffic_audit"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "spike_condition.0.spike_ratio", "100"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "spike_condition.0.spike_type", "down"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "spike_condition.0.time_diff", "5"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "window", "5"),
				),
			},
			{
				Config: notificationChannelConf + spikeConditionAlertStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_alert.alert", "id", regexp.MustCompile("^alr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "name", "alert-name1"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "description", "alert-description1"),
					resource.TestMatchResourceAttr("pfptmeta_alert.alert", "channels.0", regexp.MustCompile("^nch-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "group_by", "src_id"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "notify_message", "metaport disconnected1"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "query_text", "event:keepalive"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "source_type", "traffic_audit"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "spike_condition.0.spike_ratio", "99"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "spike_condition.0.spike_type", "down"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "spike_condition.0.time_diff", "60"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert", "window", "60"),
				),
			},
			{
				Config: notificationChannelConf + thresholdConditionAlert,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_alert.alert2", "id", regexp.MustCompile("^alr-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert2", "name", "alert-name1"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert2", "description", "alert-description1"),
					resource.TestMatchResourceAttr("pfptmeta_alert.alert2", "channels.0", regexp.MustCompile("^nch-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert2", "notify_message", "device onboarded"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert2", "query_text", "action:CREATE AND resource_type:Device"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert2", "source_type", "api_audit"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert2", "threshold_condition.0.formula", "count"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert2", "threshold_condition.0.op", "greater"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert2", "threshold_condition.0.threshold", "0"),
					resource.TestCheckResourceAttr("pfptmeta_alert.alert2", "window", "60"),
				),
			},
		},
	})
}

func TestAccDataSourceAlert(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("alert", "v1/alerts"),
		Steps: []resource.TestStep{
			{
				Config: notificationChannelConf + spikeConditionAlert + alertDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_alert.alert", "id", regexp.MustCompile("^alr-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_alert.alert", "name", "alert-name"),
					resource.TestCheckResourceAttr("data.pfptmeta_alert.alert", "description", "alert-description"),
					resource.TestMatchResourceAttr("data.pfptmeta_alert.alert", "channels.0", regexp.MustCompile("^nch-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_alert.alert", "group_by", "src_id"),
					resource.TestCheckResourceAttr("data.pfptmeta_alert.alert", "notify_message", "metaport disconnected"),
					resource.TestCheckResourceAttr("data.pfptmeta_alert.alert", "query_text", "event:keepalive AND src_type:MetaPort"),
					resource.TestCheckResourceAttr("data.pfptmeta_alert.alert", "source_type", "traffic_audit"),
					resource.TestCheckResourceAttr("data.pfptmeta_alert.alert", "spike_condition.0.spike_ratio", "100"),
					resource.TestCheckResourceAttr("data.pfptmeta_alert.alert", "spike_condition.0.spike_type", "down"),
					resource.TestCheckResourceAttr("data.pfptmeta_alert.alert", "spike_condition.0.time_diff", "5"),
					resource.TestCheckResourceAttr("data.pfptmeta_alert.alert", "window", "5"),
				),
			},
		},
	})
}
