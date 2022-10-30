package metaport_cluster_mapped_elements_attachment

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Description:   "Attaches mapped elements to metaport cluster.",
		ReadContext:   readResource,
		CreateContext: createResource,
		DeleteContext: deleteResource,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metaport_cluster_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "mpc"),
				ForceNew:         true,
			},
			"mapped_elements": {
				Description: "Mapped element IDs to be attached to the metaport cluster (Mapped Subnet, Mapped Service or Enterprise DNS)",
				Required:    true,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(true, "ne", "ed"),
				},
				MinItems: 1,
				ForceNew: true,
			},
		},
	}
}
