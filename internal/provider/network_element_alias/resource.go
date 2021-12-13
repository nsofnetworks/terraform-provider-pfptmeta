package network_element_alias

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext:   networkElementsAliasRead,
		CreateContext: networkElementAliasCreate,
		DeleteContext: networkElementAliasDelete,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_element_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(true, "ne"),
				ForceNew:         true,
			},
			"alias": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateWildcardHostName(),
				ForceNew:         true,
			},
		},
	}
}
