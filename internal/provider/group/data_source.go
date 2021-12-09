package group

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Groups represent a collection of users, typically belong to a common department or share same privileges in the organization.",

		ReadContext: groupRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"name"},
				ValidateDiagFunc: common.ValidateID(false, "grp"),
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
			"expression": {
				Description: "Allows grouping entities by their tags. Filtering by tag value is also supported if provided. " +
					"Supported operations: AND, OR, parenthesis.",
				Type:     schema.TypeString,
				Computed: true,
			},
			"provisioned_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
