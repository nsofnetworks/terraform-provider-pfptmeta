package metaport_failover

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
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
				Description: "List of mapped element IDs",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(true, "ed", "ne")},
			},
			"cluster_1": {
				Description: "Priority #1 MetaPort cluster ID. This cluster is active by default. " +
					"When failover condition is met for this cluster, the higher priority cluster becomes active.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateID(false, "mpc"),
			},
			"cluster_2": {
				Description:      "Priority #2 MetaPort cluster ID. This cluster becomes active, when failover condition is met for a lower priority cluster.",
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateID(false, "mpc"),
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
							Description:      "ENUM: [auto, manual], defaults to auto.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateENUM("auto", "manual"),
						},
					},
				},
				Optional: true,
				MaxItems: 1,
			},
			"failover": {
				Description: "Secondary to primary cluster switchover.",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"trigger": {
							Description:      "ENUM: [auto, manual], defaults to auto.",
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateENUM("auto", "manual"),
						},
						"delay": {
							Description: "Number of minutes to wait before execution of failover, defaults to 1.",
							Type:        schema.TypeInt,
							Required:    true,
						},
						"threshold": {
							Description: "Minimum number of healthy MetaPorts to keep a cluster active. Zero (0) denotes all MetaPorts in a cluster.",
							Type:        schema.TypeInt,
							Required:    true,
						},
					},
				},
				Optional: true,
				MaxItems: 1,
			},
		},
	}
}
