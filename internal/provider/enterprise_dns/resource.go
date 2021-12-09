package enterprise_dns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Enterprise DNS provides integration with global, enterprise DNS servers, " +
			"allowing resolution of FQDNs for domains that are in different locations/datacenters.",

		CreateContext: enterpriseDNSCreate,
		ReadContext:   enterpriseDNSRead,
		UpdateContext: enterpriseDNSUpdate,
		DeleteContext: enterpriseDNSDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		}, Schema: map[string]*schema.Schema{
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
			"mapped_domains": {
				Description: "DNS suffixes to be resolved within the enterprise DNS server",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mapped_domain": {
							Description:      "Proofpoint DNS Suffix",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateHostName(),
						},
						"name": {
							Description:      "Enterprise DNS server DNS suffix",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateHostName(),
						},
					},
				},
				Required: true,
			},
		},
	}
}
