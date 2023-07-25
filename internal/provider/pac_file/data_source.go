package pac_file

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: pacFileRead,
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
			"has_content": {
				Description: hasContentDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"priority": {
				Description: priorityDesc,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"content": {
				Description: contentDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"type": {
				Description: contentDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"managed_content": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domains": {
							Description: managedContentDomainsDesc,
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1000,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"cloud_apps": {
							Description: managedContentCloudAppsDesc,
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1000,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"ip_networks": {
							Description: managedContentIPNetworksDesc,
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1000,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}
