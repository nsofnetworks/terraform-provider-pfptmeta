package aac_rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: aacRuleCreate,
		ReadContext:   aacRuleRead,
		UpdateContext: aacRuleUpdate,
		DeleteContext: aacRuleDelete,
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
			"priority": {
				Description:      priorityDesc,
				Type:             schema.TypeInt,
				ValidateDiagFunc: common.ValidateIntRange(1, 5000),
				Required:         true,
			},
			"action": {
				Description:      actionDesc,
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateStringENUM("block", "allow", "isolate_block_file_transfers"),
			},
			"app_ids": {
				Description: appIdsDesc,
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "app")},
				ConflictsWith: []string{"apply_all_apps"},
			},
			"apply_all_apps": {
				Description:   applyAllAppsDesc,
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"app_ids"},
			},
			"sources": {
				Description: sourcesDesc,
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "usr", "grp")},
			},
			"exempt_sources": {
				Description: exemptSources,
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "usr", "grp"),
				},
			},
			"suspicious_login": {
				Description:      suspiciousLoginDesc,
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateStringENUM("suspicious", "safe", "any"),
			},
			"filter_expression": {
				Description: expressionDesc,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"networks": {
				Description: networksDesc,
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "ipn"),
				},
			},
			"locations": {
				Description: locationsDesc,
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"ip_reputations": {
				Description: IPDeputationsDesc,
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("tor", "data_center_hosting",
						"proxy", "vpn", "undistinguished"),
				},
			},
			"certificate_id": {
				Description:      CertificateIdDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateID(false, "crt"),
			},
			"notification_channels": {
				Description: notificationChannelsDesc,
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "nch")},
			},
		},
	}
}
