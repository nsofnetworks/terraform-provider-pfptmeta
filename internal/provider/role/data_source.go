package role

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext: roleRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"name"},
			},
			"name": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"id"},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"privileges": {
				Description: privilegesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"apply_to_orgs": {
				Description: applyToOrgsDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"all_read_privileges": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"all_write_privileges": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"all_suborgs": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"suborgs_expression": {
				Description: subOrgsExpressionDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
