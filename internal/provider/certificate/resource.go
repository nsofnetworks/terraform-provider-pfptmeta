package certificate

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

const (
	certDesc = "SSL certificate in PEM format used for BYO CA"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: certificateCreate,
		ReadContext:   certificateRead,
		UpdateContext: certificateUpdate,
		DeleteContext: certificateDelete,
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
			"certificate": {
				Description:  certDesc,
				ForceNew:     true,
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"certificate", "sans"},
			},
			"sans": {
				Description: sansDesc,
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateDNS(),
				},
				ForceNew:     true,
				MinItems:     1,
				ExactlyOneOf: []string{"certificate", "sans"},
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
