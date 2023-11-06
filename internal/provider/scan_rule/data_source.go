package scan_rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: scanRuleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"priority": {
				Description: priorityDesc,
				Type:        schema.TypeInt,
				Computed:    true,
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
			"apply_to_org": {
				Description: applyToOrgDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"content_categories": {
				Description: contentCategoriesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"threat_categories": {
				Description: threatCategoriesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"cloud_apps": {
				Description: cloudAppsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"cloud_app_risk_groups": {
				Description: cloudAppRiskGroupsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"catalog_app_categories": {
				Description: catalogAppCategoriesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"networks": {
				Description: networksDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"countries": {
				Description: countriesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"user_agents": {
				Description: userAgentsDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"user_actions": {
				Description: userActionsDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"all_supported_file_types": {
				Description: allSupportedFileTypeDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"file_types": {
				Description: fileTypesDesc,
				Type:        schema.TypeList,

				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"max_file_size_mb": {
				Description: maxFileSizeMBDesc,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"password_protected_files": {
				Description: passwordProtectedFilesDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"dlp": {
				Description: dlpDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"detectors": {
				Description: detectorsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"malware": {
				Description: malwareDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"sandbox": {
				Description: sandboxDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"antivirus": {
				Description: antivirusDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"action": {
				Description: actionDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
