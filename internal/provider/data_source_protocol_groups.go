package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceProtocolGroups() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Protocol Groups are protocols and ports that are necessary to include in granular policies",

		ReadContext: protocolGroupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protocols": {
				Description: "A list of protocols",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"from_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"to_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"proto": {
							Description: "tcp, udp or icmp",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
			"read_only": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}
