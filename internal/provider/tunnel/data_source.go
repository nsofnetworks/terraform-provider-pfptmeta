package tunnel

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		Description: description,

		ReadContext: getTunnelRead(true),
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Optional:         true,
				ExactlyOneOf:     []string{"name"},
				ValidateDiagFunc: common.ValidateID(false, "tun"),
			},
			"name": {
				ExactlyOneOf: []string{"id"},
				Type:         schema.TypeString,
				Optional:     true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"gre_config": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"source_ips": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
				Optional: true,
				MaxItems: 1,
			},
		},
	}
}
