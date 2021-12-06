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
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: validateID(false, "mpf"),
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
				Description: "List of mapped element IDs",
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"cluster_1": {
				Description: "Priority #1 MetaPort cluster ID. This cluster is active by default. " +
					"When failover condition is met for this cluster, the higher priority cluster becomes active.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_2": {
				Description: "Priority #2 MetaPort cluster ID. This cluster becomes active, when failover condition is met for a lower priority cluster.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"active_cluster": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"failback": {
				Description: "Primary to secondary cluster switchover.",
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
				Description: "Secondary to primary cluster switchover.",
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
							Description: "Minimum number of healthy MetaPorts to keep a cluster active. Zero (0) denotes all MetaPorts in a cluster.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
				Computed: true,
			},
		},
	}
}
