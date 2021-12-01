package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMappedHost() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Additional domain names for specific hosts on the mapped subnet",

		ReadContext:   mappedHostRead,
		CreateContext: mappedHostCreate,
		DeleteContext: mappedHostDelete,
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
			"mapped_host": {
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
