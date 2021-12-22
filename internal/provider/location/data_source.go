package location

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: locationRead,
		Description: "Locations are Points-of-Presence (POPs). " +
			"PoPs ensure lower latency for users while improving location awareness to make ads more relevant. " +
			"Moreover, they enhance data privacy for local users. " +
			"In addition, egress rules can be defined to route traffic via specific location.",
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"city": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"country": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"latitude": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"longitude": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
