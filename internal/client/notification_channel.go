package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"net/http"
)

const (
	notificationChannelEndpoint string = "v1/notification_channels"
)

type EmailConfig struct {
	Recipients []string `json:"recipients"`
}

func newEmailConfig(d *schema.ResourceData) *EmailConfig {
	res := &EmailConfig{}
	mailConf := d.Get("email_config").([]interface{})
	if len(mailConf) == 0 {
		return nil
	}
	r := mailConf[0].(map[string]interface{})["recipients"].([]interface{})
	res.Recipients = make([]string, len(r))
	for i, val := range r {
		res.Recipients[i] = val.(string)
	}
	return res
}

type PagerdutyConfig struct {
	ApiKey string `json:"api_key"`
}

func newPagerdutyConfig(d *schema.ResourceData) *PagerdutyConfig {
	res := &PagerdutyConfig{}
	pdConf := d.Get("pagerduty_config").([]interface{})
	if len(pdConf) == 0 {
		return nil
	}
	a := pdConf[0].(map[string]interface{})["api_key"]
	res.ApiKey = a.(string)
	return res
}

type SlackConfig struct {
	Channel string `json:"channel"`
	Url     string `json:"url"`
}

func newSlackConfig(d *schema.ResourceData) *SlackConfig {
	res := &SlackConfig{}
	sConf := d.Get("slack_config").([]interface{})
	if len(sConf) == 0 {
		return nil
	}
	conf := sConf[0].(map[string]interface{})
	res.Channel = conf["channel"].(string)
	res.Url = conf["url"].(string)
	return res
}

type Oauth2Config struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	TokenUrl     string `json:"token_url"`
}

type Auth struct {
	Oauth2Config Oauth2Config `json:"oauth2_config"`
}

type WebhookConfig struct {
	Auth          *Auth    `json:"auth,omitempty"`
	CustomPayload string   `json:"custom_payload,omitempty"`
	Headers       []string `json:"headers,omitempty"`
	Method        string   `json:"method"`
	Url           string   `json:"url"`
}

func NewWebhookConfig(d *schema.ResourceData) *WebhookConfig {
	res := &WebhookConfig{}
	whConf := d.Get("webhook_config").([]interface{})
	if len(whConf) == 0 {
		return nil
	}
	conf := whConf[0].(map[string]interface{})
	res.CustomPayload = conf["custom_payload"].(string)
	headers := conf["headers"].([]interface{})
	res.Headers = make([]string, len(headers))
	for i, val := range headers {
		res.Headers[i] = val.(string)
	}
	res.Method = conf["method"].(string)
	res.Url = conf["url"].(string)
	auth := conf["auth"].([]interface{})[0].(map[string]interface{})
	auth2Conf := auth["oauth2_config"].([]interface{})[0].(map[string]interface{})
	res.Auth = &Auth{
		Oauth2Config: Oauth2Config{
			ClientId:     auth2Conf["client_id"].(string),
			ClientSecret: auth2Conf["client_secret"].(string),
			TokenUrl:     auth2Conf["token_url"].(string),
		},
	}
	return res
}

type NotificationChannel struct {
	ID              string           `json:"id,omitempty"`
	Name            string           `json:"name,omitempty"`
	Description     string           `json:"description"`
	Enabled         bool             `json:"enabled"`
	Type            string           `json:"type,omitempty"`
	EmailConfig     *EmailConfig     `json:"email_config,omitempty"`
	PagerdutyConfig *PagerdutyConfig `json:"pagerduty_config,omitempty"`
	SlackConfig     *SlackConfig     `json:"slack_config,omitempty"`
	WebhookConfig   *WebhookConfig   `json:"webhook_config,omitempty"`
}

func NewNotificationChannel(d *schema.ResourceData) *NotificationChannel {
	res := &NotificationChannel{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Enabled = d.Get("enabled").(bool)
	res.EmailConfig = newEmailConfig(d)
	res.PagerdutyConfig = newPagerdutyConfig(d)
	res.SlackConfig = newSlackConfig(d)
	res.WebhookConfig = NewWebhookConfig(d)
	return res
}

func parseNotificationChannel(resp *http.Response) (*NotificationChannel, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	nc := &NotificationChannel{}
	err = json.Unmarshal(body, nc)
	if err != nil {
		return nil, fmt.Errorf("could not parse notification channel response: %v", err)
	}
	return nc, nil
}

func CreateNotificationChannel(ctx context.Context, c *Client, nc *NotificationChannel) (*NotificationChannel, error) {
	ncUrl := fmt.Sprintf("%s/%s", c.BaseURL, notificationChannelEndpoint)
	body, err := json.Marshal(nc)
	if err != nil {
		return nil, fmt.Errorf("could not convert notification channel to json: %v", err)
	}
	resp, err := c.Post(ctx, ncUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseNotificationChannel(resp)
}

func UpdateNotificationChannel(ctx context.Context, c *Client, ncID string, nc *NotificationChannel) (*NotificationChannel, error) {
	ncUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, notificationChannelEndpoint, ncID)
	body, err := json.Marshal(nc)
	if err != nil {
		return nil, fmt.Errorf("could not convert notification channel to json: %v", err)
	}
	resp, err := c.Patch(ctx, ncUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseNotificationChannel(resp)
}

func GetNotificationChannel(ctx context.Context, c *Client, ncID string) (*NotificationChannel, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, notificationChannelEndpoint, ncID)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseNotificationChannel(resp)
}

func DeleteNotificationChannel(ctx context.Context, c *Client, ncID string) (*NotificationChannel, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, notificationChannelEndpoint, ncID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseNotificationChannel(resp)
}
