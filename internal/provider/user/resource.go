package user

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
	"regexp"
)

var phonePattern = regexp.MustCompile("\\+[1-9]\\d{4,16}$")

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext:   userRead,
		CreateContext: userCreate,
		UpdateContext: userUpdate,
		DeleteContext: userDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"given_name": {
				Description: givenNameDesc,
				Type:        schema.TypeString,
				Required:    true,
			},
			"family_name": {
				Description: familyNameDesc,
				Type:        schema.TypeString,
				Required:    true,
			},
			"email": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateEmail(),
			},
			"phone": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidatePattern(phonePattern),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Description: tagsDesc,
				Type:        schema.TypeMap,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidatePattern(common.TagPattern)},
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}
