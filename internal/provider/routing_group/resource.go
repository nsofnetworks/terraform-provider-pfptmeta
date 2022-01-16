package routing_group

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		CreateContext: routingGroupCreate,
		ReadContext:   routingGroupRead,
		UpdateContext: routingGroupUpdate,
		DeleteContext: routingGroupDelete,
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
			"mapped_elements_ids": {
				Description: mappedElementIdsDesc,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(true, "ne"),
				},
				Optional: true,
				Computed: true,
				MinItems: 1,
			},
			"sources": {
				Description: sourcesDesc,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "ne", "usr", "grp"),
				},
				Optional: true,
				MinItems: 1,
			},
			"exempt_sources": {
				Description: exemptSourcesDesc,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "ne", "usr", "grp"),
				},
				Optional: true,
			},
		},
	}
}
