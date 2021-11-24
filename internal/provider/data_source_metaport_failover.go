package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceMetaportFailover() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "MetaPort failover defines a failover model between a primary and a secondary MetaPort clusters",

		ReadContext: metaportFailoverRead,
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
			"mapped_elements": {
				Description: "List of mapped element ids.",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"cluster_1": {
				Description: "Priority #1 metaport cluster id. This cluster is active by default. " +
					"When Failover Condition is met for this cluster the higher priority cluster becomes active.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_2": {
				Description: "Priority #2 metaport cluster id. This cluster becomes active when " +
					"failover condition is met on a lower priority cluster.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"active_cluster": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failback": {
				Description: "Primary to Secondary cluster switchover.",
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
				Description: "Secondary to Primary cluster switchover.",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delay": {
							Description: "Number of minutes to wait before execution of failover.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"threshold": {
							Description: "Minimum number of healthy metaports to keep/make a cluster active. " +
								"Zero (0) denotes all metaports in a cluster.",
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}
