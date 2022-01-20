package access_control

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext: accessControlRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "ac"),
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
			"apply_to_entities": {
				Description: applyToEntitiesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"exempt_entities": {
				Description: exemptEntitiesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"allowed_routes": {
				Description: allowedRoutesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
		},
	}
}
