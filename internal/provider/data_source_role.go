package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceRoles() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Roles define operations on the enterprise network, such as adding and removing users, adding security policies, etc.",

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
				Description: "Privileges that should be assigned to the new role. has the following form- `resource:read/write` i.e metaports:read etc.",
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"apply_to_orgs": {
				Description: "indicates which orgs this role applies at, by default will be applied to current org.",
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
				Description: "Allows grouping entities by their tags. Filtering by tag value is also supported if provided. Supported operations: AND, OR, XOR, parenthesis.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}
