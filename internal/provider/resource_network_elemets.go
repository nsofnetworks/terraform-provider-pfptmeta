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
				Description: "Key/value attributes that can be used to group elements together to Smart Groups, and placed as target or sources in Policies",
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
				Description: "Additional domain names for specific hosts in the mapped subnet",
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
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modified_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expires_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dns_name": {
				Type:     schema.TypeString,
				Computed: true,
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
			"net_id": {
				Type:     schema.TypeInt,
				Computed: true,
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
		},
	}
}
