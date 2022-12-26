package content_category

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

const maxInt = int(^uint(0) >> 1)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: contentCategoryCreate,
		ReadContext:   contentCategoryRead,
		UpdateContext: contentCategoryUpdate,
		DeleteContext: contentCategoryDelete,
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
			"confidence_level": {
				Description:      confidenceLevelDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateStringENUM("LOW", "MEDIUM", "HIGH"),
			},
			"forbid_uncategorized_urls": {
				Description: forbidUncategorizedUrlDesc,
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"types": {
				Description: typesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM(common.ContentTypes...),
				},
				Optional: true,
			},
			"urls": {
				Description: urlsDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateCustomUrlOrIPV4(),
				},
				Optional: true,
			},
		},
	}
}
