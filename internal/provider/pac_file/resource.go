package pac_file

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: pacFileCreate,
		ReadContext:   pacFileRead,
		UpdateContext: pacFileUpdate,
		DeleteContext: pacFileDelete,
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
				Optional: true,
				Default:  true,
			},
			"apply_to_org": {
				Description:   applyToOrgDesc,
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"sources", "exempt_sources"},
			},
			"sources": {
				Description: sourcesDesc,
				Type:        schema.TypeList,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "usr", "grp"),
				},
				Optional:      true,
				ConflictsWith: []string{"apply_to_org"},
			},
			"exempt_sources": {
				Description: exemptSources,
				Type:        schema.TypeList,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "usr", "grp"),
				},
				Optional:      true,
				ConflictsWith: []string{"apply_to_org"},
			},
			"has_content": {
				Description: hasContentDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"priority": {
				Description:      priorityDesc,
				Type:             schema.TypeInt,
				ValidateDiagFunc: common.ValidateIntRange(1, 5000),
				Required:         true,
			},
			"content": {
				Description: contentDesc,
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"type": {
				Description:      typeDesc,
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateStringENUM(pacTypeManaged, pacTypeBringYourOwn),
			},
			"managed_content": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domains": {
							Description: managedContentDomainsDesc,
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1000,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidatePattern(common.DomainPattern),
							},
						},
						"cloud_apps": {
							Description: managedContentCloudAppsDesc,
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1000,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidateID(false, "ca"),
							},
						},
						"ip_networks": {
							Description: managedContentIPNetworksDesc,
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1000,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidateID(false, "ipn"),
							},
						},
					},
				},
			},
		},
	}
}
