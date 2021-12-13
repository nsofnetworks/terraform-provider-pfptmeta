package enterprise_dns

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

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
				Description: mappedDomainsDesc,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mapped_domain": {
							Description:      mappedDomainDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateHostName(),
						},
						"name": {
							Description:      mdNameDescription,
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
