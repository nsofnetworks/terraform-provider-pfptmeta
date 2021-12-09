package metaport_cluster

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "MetaPort cluster defines a group of highly-available MetaPorts that are deployed together in a single data center",

		ReadContext: metaportClusterRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "mpc"),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mapped_elements": {
				Description: "List of mapped element IDs",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        schema.TypeString,
			},
			"metaports": {
				Description: "List of MetaPort IDs",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        schema.TypeString,
			},
		},
	}
}
