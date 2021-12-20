package metaport_failover

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		CreateContext: metaportFailoverCreate,
		ReadContext:   metaportFailoverRead,
		UpdateContext: metaportFailoverUpdate,
		DeleteContext: metaportFailoverDelete,
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
			"mapped_elements": {
				Description: mappedElementsDesc,
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(true, "ed", "ne")},
			},
			"cluster_1": {
				Description:      cluster1Desc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateID(false, "mpc"),
			},
			"cluster_2": {
				Description:      cluster2Desc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateID(false, "mpc"),
			},
			"active_cluster": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failback": {
				Description: failbackDesc,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger": {
							Description:      triggerDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateStringENUM("auto", "manual"),
						},
					},
				},
				Optional: true,
				MaxItems: 1,
			},
			"failover": {
				Description: failoverDesc,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger": {
							Description:      triggerDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateStringENUM("auto", "manual"),
						},
						"delay": {
							Description: failoverDelayDesc,
							Type:        schema.TypeInt,
							Required:    true,
						},
						"threshold": {
							Description: failoverThresholdDesc,
							Type:        schema.TypeInt,
							Required:    true,
						},
					},
				},
				Optional: true,
				MaxItems: 1,
			},
			"notification_channels": {
				Description: notificationChannelsDesc,
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "nch")},
			},
		},
	}
}
