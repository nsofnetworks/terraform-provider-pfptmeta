package user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext: userRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"email"},
				ValidateDiagFunc: common.ValidateID(false, "usr"),
			},
			"given_name": {
				Description: givenNameDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"family_name": {
				Description: familyNameDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Type:             schema.TypeString,
				Optional:         true,
				ConflictsWith:    []string{"id"},
				ValidateDiagFunc: common.ValidateEmail(),
			},
			"phone": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Description: tagsDesc,
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}
