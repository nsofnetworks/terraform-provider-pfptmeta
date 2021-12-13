package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

func TestAccResourceNotificationChannelMail(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("notification_channel", "v1/notification_channels"),
		Steps: []resource.TestStep{
			{
				Config: mailNotificationStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_notification_channel.mail", "id", regexp.MustCompile("^nch-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.mail", "name", "mail-channel",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.mail", "description", "mail channel description",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.mail", "email_config.0.recipients.0", "user1@example.com",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.mail", "email_config.0.recipients.1", "user2@example.com",
					),
				),
			},
			{
				Config: mailNotificationStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_notification_channel.mail", "id", regexp.MustCompile("^nch-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.mail", "name", "mail-channel",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.mail", "description", "mail channel description",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.mail", "email_config.0.recipients.0", "user1@example.com",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.mail", "email_config.0.recipients.1", "user3@example.com",
					),
				),
			},
		},
	})
}

func TestAccResourceNotificationChannelPagerDuty(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("notification_channel", "v1/notification_channels"),
		Steps: []resource.TestStep{
			{
				Config: pagerDutyNotification,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_notification_channel.pagerduty", "id", regexp.MustCompile("^nch-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.pagerduty", "name", "pagerduty-channel",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.pagerduty", "description", "pagerduty channel description",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.pagerduty", "pagerduty_config.0.api_key", "api-key",
					),
				),
			},
		},
	})
}

func TestAccResourceNotificationChannelSlack(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("notification_channel", "v1/notification_channels"),
		Steps: []resource.TestStep{
			{
				Config: slackNotification,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_notification_channel.slack", "id", regexp.MustCompile("^nch-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.slack", "name", "slack-channel",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.slack", "description", "slack channel description",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.slack", "slack_config.0.url", "https://hooks.slack.com/services/test1/test2",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.slack", "slack_config.0.channel", "#IT",
					),
				),
			},
		},
	})
}

func TestAccResourceNotificationChannelWebHook(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("notification_channel", "v1/notification_channels"),
		Steps: []resource.TestStep{
			{
				Config: webHookNotification,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_notification_channel.webhook", "id", regexp.MustCompile("^nch-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.webhook", "name", "webhook-channel",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.webhook", "description", "webhook channel description",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.webhook", "webhook_config.0.auth.0.oauth2_config.0.client_id", "client_id",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.webhook", "webhook_config.0.auth.0.oauth2_config.0.client_secret", "client_secret",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.webhook", "webhook_config.0.auth.0.oauth2_config.0.token_url", "https://token.url.com/test",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.webhook", "webhook_config.0.custom_payload", "{\"key1\": \"value1\"}",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.webhook", "webhook_config.0.headers.0", "Header1Name:Header1Value",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.webhook", "webhook_config.0.headers.1", "Header2Name:Header2Value",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.webhook", "webhook_config.0.method", "POST",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_notification_channel.webhook", "webhook_config.0.url", "https://hooks.com/test",
					),
				),
			},
		},
	})
}

func TestAccDataSourceNotificationChannel(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("notification_channel", "v1/notification_channels"),
		Steps: []resource.TestStep{
			{
				Config: mailNotificationStep1 + datasourceNotification,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_notification_channel.mail", "id", regexp.MustCompile("^nch-.*$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_notification_channel.mail", "name", "mail-channel",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_notification_channel.mail", "description", "mail channel description",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_notification_channel.mail", "email_config.0.recipients.0", "user1@example.com",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_notification_channel.mail", "email_config.0.recipients.1", "user2@example.com",
					),
				),
			},
		},
	})
}

const (
	mailNotificationStep1 = `
resource "pfptmeta_notification_channel" "mail" {
  name = "mail-channel"
  description = "mail channel description"
  email_config {
    recipients = ["user1@example.com", "user2@example.com"]
  }
}
`
	mailNotificationStep2 = `
resource "pfptmeta_notification_channel" "mail" {
  name = "mail-channel"
  description = "mail channel description"
  email_config {
    recipients = ["user1@example.com", "user3@example.com"]
  }
}
`
	pagerDutyNotification = `
resource "pfptmeta_notification_channel" "pagerduty" {
  name = "pagerduty-channel"
  description = "pagerduty channel description"
  pagerduty_config {
    api_key = "api-key"
  }
}
`
	slackNotification = `
resource "pfptmeta_notification_channel" "slack" {
  name = "slack-channel"
  description = "slack channel description"
  slack_config {
    url = "https://hooks.slack.com/services/test1/test2"
    channel = "#IT"
  }
}

`
	webHookNotification = `
resource "pfptmeta_notification_channel" "webhook" {
  name        = "webhook-channel"
  description = "webhook channel description"
  webhook_config {
    auth {
      oauth2_config {
        client_id     = "client_id"
        client_secret = "client_secret"
        token_url     = "https://token.url.com/test"
      }
    }
    custom_payload = "{\"key1\": \"value1\"}"
    headers        = ["Header1Name:Header1Value", "Header2Name:Header2Value"]
    method         = "POST"
    url            = "https://hooks.com/test"
  }
}`

	datasourceNotification = `

data "pfptmeta_notification_channel" "mail" {
  id = pfptmeta_notification_channel.mail.id
}`
)
