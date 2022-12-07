package cloud_app

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		CreateContext: cloudAppCreate,
		ReadContext:   cloudAppRead,
		UpdateContext: cloudAppUpdate,
		DeleteContext: cloudAppDelete,

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
			"app": {
				Description:      appDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateID(false, "sia"),
			},
			"tenant": {
				Description: tenantDesc,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tenant_type": {
				Description:      tenantTypeDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateStringENUM("All", "Personal", "Corporate"),
				Default:          "All",
			},
			"urls": {
				Description: urlsDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateCustomUrlOrIPV4(),
				},
				Optional: true,
				MaxItems: 20,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
