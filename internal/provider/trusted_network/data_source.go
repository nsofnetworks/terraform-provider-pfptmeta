package trusted_network

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext: trustedNetworkRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "tn"),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"apply_to_org": {
				Description: applyToOrgDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"apply_to_entities": {
				Description: applyToEntitiesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"exempt_entities": {
				Description: exemptEntitiesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"criteria": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_ip_config": {
							Description: externalIPConfigDesc,
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addresses_ranges": {
										Computed: true,
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
								},
							},
						},
						"resolved_address_config": {
							Description: resolvedAddressConfig,
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addresses_ranges": {
										Computed: true,
										Type:     schema.TypeList,
										Elem:     &schema.Schema{Type: schema.TypeString},
									},
									"hostname": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
				Computed: true,
			},
		},
	}
}
