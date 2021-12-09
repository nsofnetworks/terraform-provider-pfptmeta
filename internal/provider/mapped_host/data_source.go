package mapped_host

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
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
				ValidateDiagFunc: common.ValidateID(true, "ne"),
			},
			"mapped_host": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateHostName(),
			},
			"name": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateHostName(),
			},
		},
	}
}
