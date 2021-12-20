package alert

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
	"regexp"
)

const maxInt = int(^uint(0) >> 1)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: alertCreate,
		ReadContext:   alertRead,
		UpdateContext: alertUpdate,
		DeleteContext: alertDelete,
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
			"channels": {
				Description: channelsDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "nch"),
				},
				Required: true,
			},
			"group_by": {
				Description:      groupByDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidatePattern(regexp.MustCompile("^[\\w-]*$")),
			},
			"notify_message": {
				Description: notifyMessageDesc,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"query_text": {
				Type:     schema.TypeString,
				Required: true,
			},
			"source_type": {
				Description:      sourceTypeDesc,
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateStringENUM("security_audit", "api_audit", "traffic_audit", "webfilter_audit", "webfilter_audit"),
			},
			"spike_condition": {
				Type:          schema.TypeList,
				ConflictsWith: []string{"threshold_condition"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_hits": {
							Description:      minHitsDesc,
							Type:             schema.TypeInt,
							Optional:         true,
							Default:          0,
							ValidateDiagFunc: common.ValidateIntRange(0, maxInt),
						},
						"spike_ratio": {
							Description:      spikeRatioDesc,
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: common.ValidateIntRange(1, 100),
						},
						"spike_type": {
							Description:      spikeTypeDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateStringENUM("up", "down", "both"),
						},
						"time_diff": {
							Description:      timeDiffDesc,
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: common.ValidateIntENUM(1, 3, 5, 60, 1440, 10080),
						},
					},
				},
				Optional: true,
				MaxItems: 1,
			},
			"threshold_condition": {
				Type:          schema.TypeList,
				ConflictsWith: []string{"spike_condition"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"formula": {
							Description:      formulaDesc,
							Type:             schema.TypeString,
							Optional:         true,
							Default:          0,
							ValidateDiagFunc: common.ValidateStringENUM("count"),
						},
						"op": {
							Description:      opDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateStringENUM("greater", "greaterequals", "less", "lessequals", "equals"),
						},
						"threshold": {
							Description: thresholdDesc,
							Type:        schema.TypeInt,
							Required:    true,
						},
					},
				},
				Optional: true,
				MaxItems: 1,
			},
			"window": {
				Description:      windowDesc,
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: common.ValidateIntENUM(1, 3, 5, 10, 30, 60, 360, 1440, 2880, 10080),
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
