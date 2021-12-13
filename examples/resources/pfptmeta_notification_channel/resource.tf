resource "pfptmeta_notification_channel" "mail" {
  name        = "mail-channel"
  description = "mail channel description"
  email_config {
    recipients = ["user1@example.com", "user2@example.com"]
  }
}

resource "pfptmeta_notification_channel" "pagerduty" {
  name        = "pagerduty-channel"
  description = "pagerduty channel description"
  pagerduty_config {
    api_key = "api-key"
  }
}

resource "pfptmeta_notification_channel" "slack" {
  name        = "slack-channel"
  description = "slack channel description"
  slack_config {
    url     = "https://hooks.slack.com/services/test1/test2"
    channel = "#IT"
  }
}

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
}