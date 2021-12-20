package alert

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: alertRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
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
			"channels": {
				Description: channelsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"group_by": {
				Description: groupByDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"notify_message": {
				Description: notifyMessageDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"query_text": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"source_type": {
				Description: sourceTypeDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"spike_condition": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"min_hits": {
							Description: minHitsDesc,
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"spike_ratio": {
							Description: spikeRatioDesc,
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"spike_type": {
							Description: spikeTypeDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"time_diff": {
							Description: timeDiffDesc,
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
			"threshold_condition": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"formula": {
							Description: formulaDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"op": {
							Description: opDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"threshold": {
							Description: thresholdDesc,
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
			"window": {
				Description: windowDesc,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
