package routing_group_mapped_elements_attachment

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Description:   "Attaches mapped elements to routing group.",
		ReadContext:   readResource,
		CreateContext: createResource,
		DeleteContext: deleteResource,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"routing_group_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "rg"),
				ForceNew:         true,
			},
			"mapped_elements_ids": {
				Description: "Mapped element IDs to be attached to the routing group (Mapped Subnet or Mapped Service)",
				Required:    true,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(true, "ne"),
				},
				MinItems: 1,
				ForceNew: true,
			},
		},
	}
}
