package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"regexp"
)

var privilegesPattern = regexp.MustCompile("^[a-z_]+:(read|write)$")

func resourceRole() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Roles define operations on the enterprise network, such as adding and removing users, adding security policies, etc.",

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
				Description: "Privileges that should be assigned to the new role. has the following form- `resource:read/write` i.e metaports:read etc.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validatePattern(privilegesPattern)},
				Optional: true,
			},
			"apply_to_orgs": {
				Description: "indicates which orgs this role applies at, by default will be applied to current org.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateID(false, "org"),
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
				Description: "Allows grouping entities by their tags. Filtering by tag value is also supported if provided. Supported operations: AND, OR, XOR, parenthesis.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
}
