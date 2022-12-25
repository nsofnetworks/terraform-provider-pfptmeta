package dlp_rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: dlpRuleRead,
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
			"action": {
				Description: actionDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"apply_to_org": {
				Description: applyToOrgDesc,
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
			"alert_level": {
				Description: alertLevelDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"all_resources": {
				Description: allResourcesDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"all_supported_file_types": {
				Description: allSupportedFileTypeDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"cloud_apps": {
				Description: cloudAppsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"content_types": {
				Description: contentTypesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"detectors": {
				Description: detectorsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"file_parts": {
				Description: filePartsDesc,
				Type:        schema.TypeList,

				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"file_types": {
				Description: fileTypesDesc,
				Type:        schema.TypeList,

				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"filter_expression": {
				Description: expressionDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"priority": {
				Description: priorityDesc,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"resource_countries": {
				Description: resourceCountriesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"threat_types": {
				Description: threatTypesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"user_actions": {
				Description: userActionsDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
