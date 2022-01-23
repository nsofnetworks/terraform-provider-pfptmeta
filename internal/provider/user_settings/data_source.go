package user_settings

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: userSettingsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "as"),
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
			"apply_on_org": {
				Description: applyOnOrgDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"apply_to_entities": {
				Description: applyToEntitiesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"max_devices_per_user": {
				Description: maxDevicesPerUserDesc,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"password_expiration": {
				Description: passwordExpirationDesc,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"prohibited_os": {
				Description: prohibitedOsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"proxy_pops": {
				Description: proxyPopsDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"sso_mandatory": {
				Description: ssoMandatoryDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"mfa_required": {
				Description: mfaRequiredDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"allowed_factors": {
				Description: allowedFactorsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
		},
	}
}
