package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProtocolGroup() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Protocol Groups are protocols and ports that must be included into granular policies.",

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
				Description: "A list of protocols",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validateIntRange(0, 65535),
						},
						"to_port": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validateIntRange(0, 65535),
						},
						"proto": {
							Description:      "Protocol type, can be one of: tcp, udp, icmp",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateENUM("tcp", "udp", "icmp"),
						},
					},
				},
				Required: true,
			},
		},
	}
}
