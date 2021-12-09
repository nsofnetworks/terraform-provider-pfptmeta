package enterprise_dns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Enterprise DNS provides integration with global, enterprise DNS servers, " +
			"allowing resolution of FQDNs for domains that are in different locations/datacenters.",

		ReadContext: enterpriseDNSRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mapped_domains": {
				Description: "DNS suffixes to be resolved within the enterprise DNS server",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mapped_domain": {
							Description: "Proofpoint DNS Suffix",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "Enterprise DNS server DNS suffix",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}
