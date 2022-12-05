package tenant_restriction

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
	"regexp"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		CreateContext: trCreate,
		ReadContext:   trRead,
		UpdateContext: trUpdate,
		DeleteContext: trDelete,
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
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"google_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"microsoft_config"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_consumer_access": {
							Description: allowConsumerAccessDesc,
							Type:        schema.TypeBool,
							Required:    true,
						},
						"allow_service_accounts": {
							Description: allowServiceAccountDesc,
							Type:        schema.TypeBool,
							Required:    true,
						},
						"tenants": {
							Description: googleTenantsDesc,
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    10,
							Required:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidateDomainName(),
							},
						},
					},
				},
			},
			"microsoft_config": {
				Type:          schema.TypeList,
				MaxItems:      1,
				Optional:      true,
				ConflictsWith: []string{"google_config"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"allow_personal_microsoft_domains": {
							Description: allowPersonalDomainsDesc,
							Type:        schema.TypeBool,
							Required:    true,
						},
						"tenant_directory_id": {
							Description: TenantDirectoryIdDesc,
							Type:        schema.TypeString,
							Required:    true,
						},
						"tenants": {
							Description: microsoftTenantsDesc,
							Type:        schema.TypeList,
							MinItems:    1,
							MaxItems:    10,
							Required:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidatePattern(regexp.MustCompile("^[a-zA-Z0-9\\-._]+$")),
							},
						},
					},
				},
			},
		},
	}
}
