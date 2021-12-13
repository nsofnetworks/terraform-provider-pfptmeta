package notification_channel

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext: ncRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recipients": {
							Description: recipientsDesc,
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"pagerduty_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Description: pagerDutyApiKeyConfDesc,
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"slack_config": {
				Description: slackURLDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Description: slackChannelDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"url": {
							Description: slackURLDesc,
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"webhook_config": {
				Description: webHookConf,
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"oauth2_config": {
										Computed: true,
										Type:     schema.TypeList,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"client_id": {
													Type:      schema.TypeString,
													Computed:  true,
													Sensitive: true,
												},
												"client_secret": {
													Description: sensitiveDataWarning,
													Type:        schema.TypeString,
													Computed:    true,
													Sensitive:   true,
												},
												"token_url": {
													Type:     schema.TypeString,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"custom_payload": {
							Description: payloadDesc,
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
						"headers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type:      schema.TypeString,
								Sensitive: true,
							},
						},
						"method": {
							Description: methodDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
