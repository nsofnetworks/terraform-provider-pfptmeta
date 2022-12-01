package time_frame

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: timeFrameRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "tmf"),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"days": {
				Description: daysDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hour": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"minute": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"end_time": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hour": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"minute": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
