package notification_channel

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

const (
	description = `The system supports four channel notification types. Notifications are sent according to the logs generated by various system components.

Channel – A conduit for delivering notifications to different systems. The following notification channels are supported:

Email – Sends notifications to an email address.

Slack – Uses Slack Webhook API to send notifications to a specific Slack channel.

PagerDuty – Uses the PagerDuty API key created in your account to Send notifications to specific PagerDuty account.

Webhook – Sends notifications to any system that supports a Webhook API.`
	recipientsDesc          = "Up to 10 email addresses to which notifications will be sent."
	pagerDutyApiKeyConfDesc = "PagerDuty API key created in your PagerDuty account. " + sensitiveDataWarning
	slackURLDesc            = "Slack Webhook URL. " + sensitiveDataWarning
	slackChannelDesc        = "Slack channel."
	webHookConf             = "Used for any system that supports Webhook API"
	payloadDesc             = "A custom JSON object to be used as a Webhook alert payload."
	methodDesc              = "Enum: POST, PUT"
	sensitiveDataWarning    = "Note that this may show up in logs, and it will be stored in the state file."
)

var ExcludedKeys = []string{"id", "email_config", "pagerduty_config", "slack_config", "webhook_config"}

func notificationChannelToResource(d *schema.ResourceData, nc *client.NotificationChannel) (diags diag.Diagnostics) {
	d.SetId(nc.ID)
	err := client.MapResponseToResource(nc, d, ExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}

	switch nc.Type {
	case "email":
		if nc.EmailConfig != nil {
			emailConfig := []map[string]interface{}{
				{"recipients": nc.EmailConfig.Recipients},
			}
			err = d.Set("email_config", emailConfig)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	case "slack":
		if nc.SlackConfig != nil {
			sConfig := []map[string]interface{}{
				{"channel": nc.SlackConfig.Channel, "url": nc.SlackConfig.Url},
			}
			err = d.Set("slack_config", sConfig)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	case "pagerduty":
		if nc.PagerdutyConfig != nil {
			pdConfig := []map[string]interface{}{
				{"api_key": nc.PagerdutyConfig.ApiKey},
			}
			err = d.Set("pagerduty_config", pdConfig)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	case "webhook":
		origWebHookVody := client.NewWebhookConfig(d)
		if nc.WebhookConfig != nil {
			whConfig := []map[string]interface{}{
				{
					"custom_payload": nc.WebhookConfig.CustomPayload,
					"headers":        nc.WebhookConfig.Headers,
					"method":         nc.WebhookConfig.Method,
					"url":            nc.WebhookConfig.Url,
				},
			}
			if nc.WebhookConfig.Auth != nil {
				whConfig[0]["auth"] = []map[string]interface{}{
					{
						"oauth2_config": []map[string]interface{}{
							{
								"client_id":     nc.WebhookConfig.Auth.Oauth2Config.ClientId,
								"client_secret": origWebHookVody.Auth.Oauth2Config.ClientSecret,
								"token_url":     nc.WebhookConfig.Auth.Oauth2Config.TokenUrl,
							},
						},
					},
				}
			}
			err = d.Set("webhook_config", whConfig)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return
}

func ncRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Get("id").(string)
	nc, err := client.GetNotificationChannel(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing notification channel %s because it's gone", id)
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	return notificationChannelToResource(d, nc)
}
func ncCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewNotificationChannel(d)
	nc, err := client.CreateNotificationChannel(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return notificationChannelToResource(d, nc)
}

func ncUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewNotificationChannel(d)
	nc, err := client.UpdateNotificationChannel(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return notificationChannelToResource(d, nc)
}

func ncDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteNotificationChannel(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return
}
