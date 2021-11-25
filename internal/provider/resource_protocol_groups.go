package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProtocolGroup() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Protocol Groups are protocols and ports that are necessary to include in granular policies",

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
							Type:     schema.TypeInt,
							Required: true,
						},
						"to_port": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"proto": {
							Description:      "tcp, udp or icmp",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateENUM("tcp", "udp", "icmp"),
						},
					},
				},
				Required: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}
