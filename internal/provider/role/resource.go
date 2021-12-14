package role

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		CreateContext: roleCreate,
		ReadContext:   roleRead,
		UpdateContext: roleUpdate,
		DeleteContext: roleDelete,
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
			"privileges": {
				Description: privilegesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidatePattern(common.PrivilegesPattern)},
				Optional: true,
				Computed: true,
			},
			"apply_to_orgs": {
				Description: applyToOrgsDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "org"),
				},
				Optional: true,
			},
			"all_read_privileges": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"all_write_privileges": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"all_suborgs": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"suborgs_expression": {
				Description: subOrgsExpressionDesc,
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}
