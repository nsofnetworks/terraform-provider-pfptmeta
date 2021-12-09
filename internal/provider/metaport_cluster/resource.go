package metaport_cluster

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "MetaPort cluster defines a group of highly-available MetaPorts that are deployed together in a single data center",

		CreateContext: metaportClusterCreate,
		ReadContext:   metaportClusterRead,
		UpdateContext: metaportClusterUpdate,
		DeleteContext: metaportClusterDelete,
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
			"mapped_elements": {
				Description: "List of mapped element IDs",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(true, "ed", "ne")},
			},
			"metaports": {
				Description: "List of MetaPort IDs",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(true, "mp")},
			},
		},
	}
}
