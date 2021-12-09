package group_roles_attachment

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		ReadContext:   readResource,
		CreateContext: createResource,
		DeleteContext: deleteResource,
		UpdateContext: createResource,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "grp"),
				ForceNew:         true,
			},
			"roles": {
				Description: "Role IDs that will be attached to the group",
				Required:    true,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "rol"),
				},
				MinItems: 1,
			},
		},
	}
}
