package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMappedHost() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Additional domain names for specific hosts on the mapped subnet",

		ReadContext: mappedHostRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_element_id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateID(true, "ne"),
			},
			"mapped_host": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateHostName(),
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateHostName(),
			},
		},
	}
}
