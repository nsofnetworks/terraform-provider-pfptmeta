package network_element

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		CreateContext: networkElementCreate,
		ReadContext:   networkElementsRead,
		UpdateContext: networkElementUpdate,
		DeleteContext: networkElementDelete,
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
			"tags": {
				Description: tagsDesc,
				Type:        schema.TypeMap,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidatePattern(common.TagPattern)},
				Optional: true,
			},
			"mapped_subnets": {
				Description: mappedSubnetsDesc,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateCIDR4(),
				},
				Optional:      true,
				ConflictsWith: []string{"mapped_service", "platform", "owner_id"},
			},
			"mapped_service": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"mapped_subnets", "platform", "owner_id"},
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Description: enabledDesc,
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"auto_aliases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"platform": {
				Description:      "One of ['Android', 'macOS', 'iOS', 'Linux', 'Windows', 'ChromeOS', 'Unknown']",
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateDiagFunc: common.ValidateStringENUM("Android", "macOS", "iOS", "Linux", "Windows", "ChromeOS", "Unknown"),
				ForceNew:         true,
			},
			"owner_id": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"mapped_subnets", "mapped_service"},
				ForceNew:         true,
				ValidateDiagFunc: common.ValidateID(false, "usr"),
			},
		},
	}
}
