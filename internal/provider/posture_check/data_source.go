package posture_check

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext: postureCheckRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "pc"),
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
			"check": {
				Description: checkDesc,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_version": {
							Description: minVersionDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"when": {
				Description: whenDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
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
			"osquery": {
				Description: osQueryDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"platform": {
				Description: platformDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_message_on_fail": {
				Description: userMessageDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"interval": {
				Description: intervalDesc,
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}
