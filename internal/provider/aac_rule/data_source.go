package aac_rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: aacRuleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "arl"),
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
			"priority": {
				Description: priorityDesc,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"action": {
				Description: actionDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"app_ids": {
				Description: appIdsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"apply_all_apps": {
				Description: applyAllAppsDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"sources": {
				Description: sourcesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"exempt_sources": {
				Description: exemptSources,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"filter_expression": {
				Description: expressionDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"networks": {
				Description: networksDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"locations": {
				Description: locationsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"ip_reputations": {
				Description: IPDeputationsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"certificate_id": {
				Description: CertificateIdDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"notification_channels": {
				Description: notificationChannelsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
		},
	}
}
