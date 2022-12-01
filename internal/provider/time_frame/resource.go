package time_frame

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: timeFrameCreate,
		ReadContext:   timeFrameRead,
		UpdateContext: timeFrameUpdate,
		DeleteContext: timeFrameDelete,
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
			"days": {
				Description: daysDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("sunday", "monday", "tuesday", "wednesday",
						"thursday", "friday", "saturday"),
				},
				MaxItems: 7,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hour": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: common.ValidateIntRange(0, 23),
						},
						"minute": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: common.ValidateIntRange(0, 59),
						},
					},
				},
			},
			"end_time": {
				Description: endTimeDesc,
				Type:        schema.TypeList,
				MaxItems:    1,
				Required:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hour": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: common.ValidateIntRange(0, 23),
						},
						"minute": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: common.ValidateIntRange(0, 59),
						},
					},
				},
			},
		},
	}
}
