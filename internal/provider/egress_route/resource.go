package egress_route

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		CreateContext: egressRouteCreate,
		ReadContext:   egressRouteRead,
		UpdateContext: egressRouteUpdate,
		DeleteContext: egressRouteDelete,

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
			"enabled": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"destinations": {
				Description: destinationsDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ComposeOrValidations(common.ValidateHostName(), common.ValidatePattern(regexp.MustCompile("^\\.$"))),
				},
				Optional: true,
			},
			"exempt_sources": {
				Description: exemptSourcesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "usr", "grp", "ne", "dev", "mc"),
				},
				Optional: true,
			},
			"sources": {
				Description: sourcesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "usr", "grp", "ne", "dev", "mc"),
				},
				Optional: true,
			},
			"via": {
				Description: viaDesc,
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
}
