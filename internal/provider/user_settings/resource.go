package user_settings

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

const maxInt = int(^uint(0) >> 1)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: userSettingsCreate,
		ReadContext:   userSettingsRead,
		UpdateContext: userSettingsUpdate,
		DeleteContext: userSettingsDelete,
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
				MinItems: 1,
				Optional: true,
			},
			"max_devices_per_user": {
				Description:      maxDevicesPerUserDesc,
				Type:             schema.TypeString,
				ValidateDiagFunc: common.ValidateStringToIntRange(0, maxInt),
				Optional:         true,
				ForceNew:         true,
			},
			"prohibited_os": {
				Description: prohibitedOsDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("Android", "macOS", "iOS", "Linux", "Windows", "ChromeOS"),
				},
				Optional: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 6,
			},
			"proxy_pops": {
				Description:      proxyPopsDesc,
				Type:             schema.TypeString,
				ValidateDiagFunc: common.ValidateStringENUM("ALL_POPS", "POPS_WITH_DEDICATED_IPS"),
				Optional:         true,
				ForceNew:         true,
			},
			"password_expiration": {
				Description:      passwordExpirationDesc,
				Type:             schema.TypeInt,
				ValidateDiagFunc: common.ValidateIntRange(1, maxInt),
				Optional:         true,
				ForceNew:         true,
			},
			"sso_mandatory": {
				Description: ssoMandatoryDesc,
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"mfa_required": {
				Description: mfaRequiredDesc,
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"allowed_factors": {
				Description: allowedFactorsDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("SMS", "SOFTWARE_TOTP", "VOICECALL", "EMAIL"),
				},
				Optional: true,
				ForceNew: true,
				MinItems: 1,
				MaxItems: 4,
			},
		},
	}
}
