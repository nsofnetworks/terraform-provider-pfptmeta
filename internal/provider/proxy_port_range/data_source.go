package proxy_port_range

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: proxyPortRangeRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "ppr"),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"proto": {
				Description: proto,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"from_port": {
				Description: from_port,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"to_port": {
				Description: to_port,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}
