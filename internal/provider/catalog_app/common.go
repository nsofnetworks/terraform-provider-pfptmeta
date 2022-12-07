package catalog_app

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
)

var excludedKeys = []string{"attributes"}

func catalogAppRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	name := d.Get("name").(string)
	category := d.Get("category").(string)
	ca, err := client.GetCatalogAppByName(ctx, c, name, category)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ca.ID)
	err = client.MapResponseToResource(ca, d, excludedKeys)
	attributes := []map[string]interface{}{
		{"tenant_awareness_data": []map[string]bool{{
			"tenant_corp_id_support": ca.Attributes.TenantAwarenessData.TenantCorpIdSupport,
			"tenant_type_support":    ca.Attributes.TenantAwarenessData.TenantTypeSupport,
		}}},
	}
	err = d.Set("attributes", attributes)
	if err != nil {
		return diag.FromErr(err)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	return
}
