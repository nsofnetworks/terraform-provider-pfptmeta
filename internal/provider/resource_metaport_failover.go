package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMetaportFailover() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "MetaPort failover defines a failover model between a primary and a secondary MetaPort clusters.",

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
				Description: "List of mapped element ids.",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateID(true, "ed", "ne")},
			},
			"cluster_1": {
				Description: "Priority #1 metaport cluster id. This cluster is active by default. " +
					"When Failover Condition is met for this cluster the higher priority cluster becomes active.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateID(false, "mpc"),
			},
			"cluster_2": {
				Description: "Priority #2 metaport cluster id. This cluster becomes active when " +
					"failover condition is met on a lower priority cluster.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: validateID(false, "mpc"),
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
							Description:      "ENUM: [auto, manual], defaults to auto.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateENUM("auto", "manual"),
						},
					},
				},
				Optional: true,
				MaxItems: 1,
			},
			"failover": {
				Description: "Secondary to Primary cluster switchover.",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger": {
							Description:      "ENUM: [auto, manual], defaults to auto.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validateENUM("auto", "manual"),
						},
						"delay": {
							Description: "Number of minutes to wait before execution of failover, defaults to 1.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"threshold": {
							Description: "Minimum number of healthy metaports to keep/make a cluster active. " +
								"Zero (0) denotes all metaports in a cluster.",
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
				Optional: true,
				MaxItems: 1,
			},
		},
	}
}
