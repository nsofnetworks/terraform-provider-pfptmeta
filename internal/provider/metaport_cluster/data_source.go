package metaport_cluster

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext: metaportClusterDataSourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				ExactlyOneOf:     []string{"name"},
				Optional:         true,
				ValidateDiagFunc: common.ValidateID(false, "mpc"),
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
			"mapped_elements": {
				Description: mappedElementsDesc,
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"metaports": {
				Description: metaportsDesc,
				Type:        schema.TypeSet,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
