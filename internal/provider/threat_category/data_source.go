package threat_category

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: threatCategoryRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "tc"),
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
			"risk_level": {
				Description: riskLevelDesc,
				Type:        schema.TypeString,
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
			"countries": {
				Description: countriesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"third_party_app": {
				Description: thirdPartyAppDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
