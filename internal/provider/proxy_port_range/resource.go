package proxy_port_range

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: proxyPortRangeCreate,
		ReadContext:   proxyPortRangeRead,
		UpdateContext: proxyPortRangeUpdate,
		DeleteContext: proxyPortRangeDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"proto": {
				Description:      proto,
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateStringENUM("HTTP", "HTTPS"),
			},
			"from_port": {
				Description:      from_port,
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: common.ValidateIntRange(1, 65535),
			},
			"to_port": {
				Description:      to_port,
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: common.ValidateIntRange(1, 65535),
			},
			"read_only": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}
