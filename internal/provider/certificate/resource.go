package certificate

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: certificateCreate,
		ReadContext:   certificateRead,
		UpdateContext: certificateUpdate,
		DeleteContext: certificateDelete,
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
			"sans": {
				Description: sansDesc,
				Type:        schema.TypeSet,
				Required:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateDomainName(),
				},
				ForceNew: true,
				MinItems: 1,
			},
			"serial_number": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Description: stateDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"status_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_not_after": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"valid_not_before": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
