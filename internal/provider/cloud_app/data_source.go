package cloud_app

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext: cloudAppRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "ca"),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"app": {
				Description: appDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tenant": {
				Description: tenantDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"tenant_type": {
				Description: tenantTypeDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"urls": {
				Description: urlsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
