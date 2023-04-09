package device_settings

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

const maxInt = int(^uint(0) >> 1)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: deviceSettingsCreate,
		ReadContext:   deviceSettingsRead,
		UpdateContext: deviceSettingsUpdate,
		DeleteContext: deviceSettingsDelete,
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
				Optional: true,
				Default:  true,
			},
			"apply_on_org": {
				Description: applyOnOrgDesc,
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
			"auto_fqdn_domain_names": {
				Description: autoFqdnDomainNamesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateDomainName(),
				},
				MaxItems: 1,
				MinItems: 1,
				Optional: true,
				ForceNew: true,
			},
			"direct_sso": {
				Description:      directSsoDesc,
				Type:             schema.TypeString,
				ValidateDiagFunc: common.ComposeOrValidations(common.ValidateID(false, "idp"), common.ValidatePattern(regexp.MustCompile("^false$"))),
				Optional:         true,
				Computed:         true,
				ForceNew:         true,
			},
			"overlay_mfa_refresh_period": {
				Description:      overlayMFARefreshPeriodDesc,
				Type:             schema.TypeInt,
				ValidateDiagFunc: common.ValidateIntRange(10, maxInt),
				Optional:         true,
				ForceNew:         true,
			},
			"overlay_mfa_required": {
				Description: overlayMFARequiredDesc,
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"protocol_selection_lifetime": {
				Description:      protocolSelectionLifeTimeDesc,
				Type:             schema.TypeString,
				ValidateDiagFunc: common.ValidateStringToIntRange(0, 525600),
				Optional:         true,
				ForceNew:         true,
			},
			"proxy_always_on": {
				Description: proxyAlwaysOnDesc,
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"search_domains": {
				Description: searchDomainsDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateHostName(),
				},
				Optional: true,
				ForceNew: true,
			},
			"session_lifetime": {
				Description:      sessionLifeTimeDesc,
				Type:             schema.TypeInt,
				ValidateDiagFunc: common.ValidateIntRange(1, maxInt),
				Optional:         true,
				ForceNew:         true,
			},
			"session_lifetime_grace": {
				Description:      sessionLifeTimeGraceDesc,
				Type:             schema.TypeString,
				ValidateDiagFunc: common.ValidateStringToIntRange(0, 60),
				Optional:         true,
				ForceNew:         true,
			},
			"tunnel_mode": {
				Description:      tunnelModeDesc,
				Type:             schema.TypeString,
				ValidateDiagFunc: common.ValidateStringENUM("split", "full"),
				Optional:         true,
				ForceNew:         true,
			},
			"vpn_login_browser": {
				Description:      vpnLoginBrowserDesc,
				Type:             schema.TypeString,
				ValidateDiagFunc: common.ValidateStringENUM("AGENT", "EXTERNAL", "USER_DEFINED"),
				Optional:         true,
				ForceNew:         true,
			},
			"ztna_always_on": {
				Description: ztnaAlwaysOnDesc,
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
		},
	}
}
