package catalog_app

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func DataSource() *schema.Resource {
	return &schema.Resource{
		ReadContext: catalogAppRead,
		Description: `To improve internet access management, organizations can use cloud applications. 
Each cloud application object is based on the Proofpoint shadow IT catalog with up 50,000 entries.

~> **NOTE:** Catalog App names aren't unique. If more than one result matches the name and category filters, the first one is selected.`,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "Name of the catalog application",
				Type:        schema.TypeString,
				Required:    true,
			},
			"category": {
				Description: "When used, the catalog application is filtered by its value as well",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			"risk": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"urls": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"vendor": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"verified": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"attributes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"tenant_awareness_data": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tenant_corp_id_support": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"tenant_type_support": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
