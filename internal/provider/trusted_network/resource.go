package trusted_network

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		CreateContext: trustedNetworkCreate,
		ReadContext:   trustedNetworkRead,
		UpdateContext: trustedNetworkUpdate,
		DeleteContext: trustedNetworkDelete,

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
			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"apply_to_org": {
				Description: applyToOrgDesc,
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"apply_to_entities": {
				Description: applyToEntitiesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ComposeOrValidations(common.ValidateID(true, "ne"), common.ValidateID(false, "usr", "grp")),
				},
				Optional: true,
			},
			"exempt_entities": {
				Description: exemptEntitiesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ComposeOrValidations(common.ValidateID(true, "ne"), common.ValidateID(false, "usr", "grp")),
				},
				Optional: true,
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
							Optional:    true,
							MinItems:    1,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addresses_ranges": {
										Type:     schema.TypeList,
										MinItems: 1,
										Required: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: common.ValidateCIDR4(),
										},
									},
								},
							},
						},
						"resolved_address_config": {
							Description: resolvedAddressConfig,
							Type:        schema.TypeList,
							Optional:    true,
							MinItems:    1,
							MaxItems:    1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"addresses_ranges": {
										Type:     schema.TypeList,
										MinItems: 1,
										Required: true,
										Elem: &schema.Schema{
											Type:             schema.TypeString,
											ValidateDiagFunc: common.ValidateCIDR4(),
										},
									},
									"hostname": {
										Type:             schema.TypeString,
										ValidateDiagFunc: common.ValidateHostName(),
										Required:         true,
									},
								},
							},
						},
					},
				},
				Optional: true,
				MinItems: 1,
			},
		},
	}
}
