package notification_channel

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		CreateContext: ncCreate,
		ReadContext:   ncRead,
		UpdateContext: ncUpdate,
		DeleteContext: ncDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"email_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"pagerduty_config", "slack_config", "webhook_config"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"recipients": {
							Description: recipientsDesc,
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    10,
							Required:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidateEmail(),
							},
						},
					},
				},
			},
			"pagerduty_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"email_config", "slack_config", "webhook_config"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key": {
							Description: pagerDutyApiKeyConfDesc,
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
						},
					},
				},
			},
			"slack_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"email_config", "pagerduty_config", "webhook_config"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"channel": {
							Description: slackChannelDesc,
							Type:        schema.TypeString,
							Optional:    true,
						},
						"url": {
							Description:      slackURLDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateURL(),
							Sensitive:        true,
						},
					},
				},
			},
			"webhook_config": {
				Description:   webHookConf,
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"email_config", "pagerduty_config", "slack_config"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auth": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"oauth2_config": {
										Required: true,
										Type:     schema.TypeList,
										MinItems: 1,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"client_id": {
													Type:      schema.TypeString,
													Required:  true,
													Sensitive: true,
												},
												"client_secret": {
													Description: sensitiveDataWarning,
													Type:        schema.TypeString,
													Required:    true,
													Sensitive:   true,
												},
												"token_url": {
													Type:             schema.TypeString,
													Required:         true,
													ValidateDiagFunc: common.ValidateURL(),
												},
											},
										},
									},
								},
							},
						},
						"custom_payload": {
							Description:      payloadDesc,
							Type:             schema.TypeString,
							Optional:         true,
							Sensitive:        true,
							ValidateDiagFunc: common.ValidateJson(),
						},
						"headers": {
							Type:     schema.TypeList,
							MinItems: 1,
							Optional: true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidatePattern(common.HttpHeaderPattern),
								Sensitive:        true,
							},
						},
						"method": {
							Description:      methodDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateStringENUM("POST", "PUT"),
						},
						"url": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateURL(),
						},
					},
				},
			},
		},
	}
}
