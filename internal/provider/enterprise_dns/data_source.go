package enterprise_dns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

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
				Description: mappedDomainsDesc,
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mapped_domain": {
							Description: mappedDomainDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: mdNameDescription,
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
