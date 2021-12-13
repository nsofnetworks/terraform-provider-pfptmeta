package network_element

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Mapped subnets are subnets available to the users within the local network, " +
			"residing behind the MetaPort. When you create a mapped subnet, you define a CIDR" +
			" and attach the subnet to a MetaPort. Optionally, you can define a dedicated host, residing on the subnet.",

		ReadContext: networkElementsRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Description: "Key/value attributes to be used for combining elements together into Smart Groups, and placed as targets or sources in Policies",
				Type:        schema.TypeMap,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"mapped_subnets": {
				Description: "CIDRs that will be mapped to the subnet",
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"mapped_service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Description: "Not allowed for mapped service and mapped domain",
				Type:        schema.TypeBool,
				Computed:    true,
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
