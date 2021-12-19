package location

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
)

func locationRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	name := d.Get("name").(string)
	l, err := client.GetLocation(ctx, c, name)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(name)
	err = client.MapResponseToResource(l, d, []string{})
	if err != nil {
		return diag.FromErr(err)
	}
	return
}
