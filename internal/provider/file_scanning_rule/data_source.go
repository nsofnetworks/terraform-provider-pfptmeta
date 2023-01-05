package file_scanning_rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: fileScanningRuleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "fsr"),
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
			"cloud_apps": {
				Description: cloudAppsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"block_all_file_types": {
				Description: blockAllFileTypesDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"block_content_types": {
				Description: blockContentTypesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"block_countries": {
				Description: blockCountriesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"block_file_types": {
				Description: blockFileTypeDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"block_threat_types": {
				Description: blockThreatTypesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"block_unsupported_files": {
				Description: blockUnsupportedFilesDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"filter_expression": {
				Description: expressionDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"malware": {
				Description: malwareDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"max_file_size_mb": {
				Description: maxFileSizeMBDesc,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"priority": {
				Description: priorityDesc,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"sandbox_file_types": {
				Description: sandboxFileTypesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"timeout_policy": {
				Description: timeoutPolicyDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
