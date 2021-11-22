package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceMetaport() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "MetaPort is a lightweight virtual appliance that enables the secure authenticated interface " +
			"interact between existing servers and the Proofpoint NaaS cloud. " +
			"Once configured, metaports enable users to access your applications via the Proofpoint cloud.",

		CreateContext: metaportCreate,
		ReadContext:   metaportRead,
		UpdateContext: metaportUpdate,
		DeleteContext: metaportDelete,
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
			"mapped_elements": {
				Description: "List of mapped element ids",
				Type:        schema.TypeSet,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: validateID(true, "ed", "ne")},
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"health": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"allow_support": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}
