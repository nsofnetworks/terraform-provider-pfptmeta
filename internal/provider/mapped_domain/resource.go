package mapped_domain

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext:   mappedDomainRead,
		CreateContext: mappedDomainCreate,
		DeleteContext: mappedDomainDelete,
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
			"mapped_domain": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateHostName(),
				ForceNew:         true,
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateHostName(),
				ForceNew:         true,
			},
		},
	}
}
