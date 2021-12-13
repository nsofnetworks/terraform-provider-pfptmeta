package protocol_group

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		CreateContext: protocolGroupCreate,
		ReadContext:   protocolGroupRead,
		UpdateContext: protocolGroupUpdate,
		DeleteContext: protocolGroupDelete,
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
			"protocols": {
				Description: protocolsDesc,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: common.ValidateIntRange(0, 65535),
						},
						"to_port": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: common.ValidateIntRange(0, 65535),
						},
						"proto": {
							Description:      protoDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateENUM("tcp", "udp", "icmp"),
						},
					},
				},
				Required: true,
			},
		},
	}
}
