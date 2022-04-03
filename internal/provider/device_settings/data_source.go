package device_settings

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: deviceSettingsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "ds"),
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
			"auto_fqdn_domain_names": {
				Description: autoFqdnDomainNamesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"direct_sso": {
				Description: directSsoDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"overlay_mfa_refresh_period": {
				Description: overlayMFARefreshPeriodDesc,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"overlay_mfa_required": {
				Description: overlayMFARequiredDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"protocol_selection_lifetime": {
				Description: protocolSelectionLifeTimeDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"proxy_always_on": {
				Description: proxyAlwaysOnDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"search_domains": {
				Description: searchDomainsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"session_lifetime": {
				Description: sessionLifeTimeDesc,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"session_lifetime_grace": {
				Description: sessionLifeTimeGraceDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tunnel_mode": {
				Description: tunnelModeDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"vpn_login_browser": {
				Description: vpnLoginBrowserDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"ztna_always_on": {
				Description: ztnaAlwaysOnDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
		},
	}
}
