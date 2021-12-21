package certificate

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: certificateRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "crt"),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sans": {
				Description: sansDesc,
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
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
