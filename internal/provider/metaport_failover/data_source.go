package metaport_failover

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext: metaportFailoverRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "mpf"),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mapped_elements": {
				Description: mappedElementsDesc,
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"cluster_1": {
				Description: cluster1Desc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"cluster_2": {
				Description: cluster2Desc,
				Type:        schema.TypeString,
				Computed:    true,
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
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
			"failover": {
				Description: failoverDesc,
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay": {
							Description: failoverDelayDesc,
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"threshold": {
							Description: failoverThresholdDesc,
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
			"notification_channels": {
				Description: notificationChannelsDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
