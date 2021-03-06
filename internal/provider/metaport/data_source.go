package metaport

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext: metaportRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Optional:         true,
				ExactlyOneOf:     []string{"name"},
				ValidateDiagFunc: common.ValidateID(true, "mp"),
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
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"allow_support": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"notification_channels": {
				Description: notificationChannelsDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}
