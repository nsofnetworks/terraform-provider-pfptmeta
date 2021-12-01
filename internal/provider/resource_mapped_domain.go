package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMappedDomain() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "DNS suffixes to be resolved within Mapped Subnet",

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
				ValidateDiagFunc: validateID(true, "ne"),
				ForceNew:         true,
			},
			"mapped_domain": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateHostName(),
				ForceNew:         true,
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateHostName(),
				ForceNew:         true,
			},
		},
	}
}
