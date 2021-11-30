package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNetworkElement() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Mapped subnets are subnets available to the users within the local network, " +
			"residing behind the MetaPort. When you create a mapped subnet, you define a CIDR" +
			" and attach the subnet to a MetaPort. Optionally, you can define a dedicated host, residing on the subnet.",

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
				Description: "Key/value attributes to be used for combining elements together into Smart Groups, and placed as targets or sources in Policies",
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"mapped_domains": {
				Description: "DNS suffixes to be resolved within this Mapped Subnet",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mapped_domain": {
							Description: "Meta DNS suffix",
							Type:        schema.TypeString,
							Required:    true,
						},
						"name": {
							Description: "Internal DNS suffix",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
				Optional: true,
				Default:  nil,
			},
			"mapped_hosts": {
				Description: "Additional domain names for specific hosts on the mapped subnet",
				Type:        schema.TypeSet,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mapped_host": {
							Description: "Hostname",
							Type:        schema.TypeString,
							Required:    true,
						},
						"name": {
							Description: "Remote hostname or IP",
							Type:        schema.TypeString,
							Required:    true,
						},
					},
				},
				Optional: true,
			},
			"mapped_subnets": {
				Description: "CIDRs that will be mapped to the subnet",
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
			},
			"mapped_service": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Description: "Not allowed for mapped service and mapped domain",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"aliases": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
				ValidateDiagFunc: validateENUM("Android", "macOS", "iOS", "Linux", "Windows", "ChromeOS", "Unknown"),
			},
			"owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
