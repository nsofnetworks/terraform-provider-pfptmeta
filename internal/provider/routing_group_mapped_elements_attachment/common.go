package routing_group_mapped_elements_attachment

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

func generateID(rgID string, meIDs []string) string {
	hash := 0
	for _, meID := range meIDs {
		hash += schema.HashString(meID)
	}
	return fmt.Sprintf("%s-%d", rgID, hash)
}

func attachmentToResource(d *schema.ResourceData, rg *client.RoutingGroup) (diags diag.Diagnostics) {
	err := d.Set("routing_group_id", rg.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	rgMes := &schema.Set{F: schema.HashString}
	for _, i := range rg.MappedElementsIds {
		rgMes.Add(i)
	}
	schemaMes := d.Get("mapped_elements_ids").(*schema.Set)
	u := schema.NewSet(schema.HashString, schemaMes.List())
	intersection := rgMes.Intersection(u)
	mes := client.ResourceTypeSetToStringSlice(intersection)
	err = d.Set("mapped_elements_ids", mes)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(generateID(rg.ID, mes))
	return
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	rgID := d.Get("routing_group_id").(string)
	rg, err := client.GetRoutingGroup(ctx, c, rgID)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	return attachmentToResource(d, rg)
}
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	rgID := d.Get("routing_group_id").(string)
	mes := client.ResourceTypeSetToStringSlice(d.Get("mapped_elements_ids").(*schema.Set))
	rg, err := client.AddMappedElementsToRoutingGroups(ctx, c, rgID, mes)
	if err != nil {
		return diag.FromErr(err)
	}
	return attachmentToResource(d, rg)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	gID := d.Get("routing_group_id").(string)
	mes := d.Get("mapped_elements_ids").(*schema.Set)
	_, err := client.RemoveMappedElementsFromRoutingGroups(ctx, c, gID, client.ResourceTypeSetToStringSlice(mes))
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId("")
	return
}
