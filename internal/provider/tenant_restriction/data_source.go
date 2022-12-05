package tenant_restriction

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext: trRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "tr"),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"google_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_consumer_access": {
							Description: allowConsumerAccessDesc,
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"allow_service_accounts": {
							Description: allowServiceAccountDesc,
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"tenants": {
							Description: googleTenantsDesc,
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidateDomainName(),
							},
						},
					},
				},
			},
			"microsoft_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_personal_microsoft_domains": {
							Description: allowPersonalDomainsDesc,
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"tenant_directory_id": {
							Description: TenantDirectoryIdDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"tenants": {
							Description: microsoftTenantsDesc,
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}
