package user_roles_attachment

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		Description:   "Attaches roles to user.",
		ReadContext:   readResource,
		CreateContext: createResource,
		DeleteContext: deleteResource,
		UpdateContext: createResource,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "usr"),
				ForceNew:         true,
			},
			"roles": {
				Description: "Role IDs to be attached to the user",
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
