package posture_check

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		CreateContext: postureCheckCreate,
		ReadContext:   postureCheckRead,
		UpdateContext: postureCheckUpdate,
		DeleteContext: postureCheckDelete,
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
				Description: enabledDesc,
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
			},
			"action": {
				Description:      actionDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateStringENUM("DISCONNECT", "NONE", "WARNING"),
				Default:          "DISCONNECT",
			},
			"check": {
				Description:  checkDesc,
				ExactlyOneOf: []string{"osquery"},
				Type:         schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_version": {
							Description:      minVersionDesc,
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: common.ValidatePattern(regexp.MustCompile("^(?P<major>0|[1-9]\\d{0,3})\\.(?P<minor>0|[1-9]\\d{0,3})\\.(?P<patch>0|[1-9]\\d{0,3})$")),
						},
						"type": {
							Description:      typeDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateStringENUM("jailbroken_rooted", "screen_lock_enabled", "minimum_app_version", "minimum_os_version", "malicious_app_detection", "developer_mode_enabled"),
						},
					},
				},
				Optional: true,
				MaxItems: 1,
			},
			"when": {
				Description: whenDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("PRE_CONNECT", "PERIODIC"),
				},
				Required: true,
			},
			"apply_to_org": {
				Description: applyToOrgDesc,
				Type:        schema.TypeBool,
				Default:     true,
				Optional:    true,
			},
			"apply_to_entities": {
				Description: applyToEntitiesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ComposeOrValidations(common.ValidateID(true, "ne", "dev"), common.ValidateID(false, "usr", "grp")),
				},
				Optional: true,
			},
			"exempt_entities": {
				Description: exemptEntitiesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ComposeOrValidations(common.ValidateID(true, "ne", "dev"), common.ValidateID(false, "usr", "grp")),
				},
				Optional: true,
			},
			"osquery": {
				Description:  osQueryDesc,
				ExactlyOneOf: []string{"check"},
				Type:         schema.TypeString,
				Optional:     true,
			},
			"platform": {
				Description:      platformDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateStringENUM("Android", "macOS", "iOS", "Linux", "Windows", "ChromeOS"),
			},
			"user_message_on_fail": {
				Description: userMessageDesc,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"interval": {
				Description:      intervalDesc,
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: common.ValidateIntENUM(0, 5, 60),
			},
		},
	}
}
