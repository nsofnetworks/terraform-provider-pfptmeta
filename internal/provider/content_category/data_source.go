package content_category

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: contentCategoryRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "cc"),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"confidence_level": {
				Description: confidenceLevelDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"forbid_uncategorized_urls": {
				Description: forbidUncategorizedUrlDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"types": {
				Description: typesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"urls": {
				Description: urlsDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ComposeOrValidations(common.ValidateURL(), common.ValidateHostName()),
				},
				Computed: true,
			},
		},
	}
}
