package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkElementAlias() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "DNS alias (FQDN) of the network element. valid for network element of type Device, Native Service and Mapped Service.",

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
				ValidateDiagFunc: validateID(true, "ne"),
				ForceNew:         true,
			},
			"alias": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateWildcardHostName(),
				ForceNew:         true,
			},
		},
	}
}
