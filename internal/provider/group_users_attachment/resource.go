package group_users_attachment

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Description:   "Adds users to group.",
		ReadContext:   readResource,
		CreateContext: createResource,
		DeleteContext: deleteResource,
		UpdateContext: updateResource,
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
			"users": {
				Description: "User IDs to be added to the group",
				Required:    true,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "usr"),
				},
				MinItems: 1,
			},
		},
	}
}
