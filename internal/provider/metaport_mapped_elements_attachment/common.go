package metaport_mapped_elements_attachment

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

func generateID(mID string, meIDs []string) string {
	hash := 0
	for _, meID := range meIDs {
		hash += schema.HashString(meID)
	}
	return fmt.Sprintf("%s-%d", mID, hash)
}

func attachmentToResource(d *schema.ResourceData, m *client.Metaport) (diags diag.Diagnostics) {
	err := d.Set("metaport_id", m.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	mMes := &schema.Set{F: schema.HashString}
	for _, i := range m.MappedElements {
		mMes.Add(i)
	}
	schemaMes := d.Get("mapped_elements").(*schema.Set)
	schemaMesSet := schema.NewSet(schema.HashString, schemaMes.List())
	intersection := mMes.Intersection(schemaMesSet)
	mes := client.ResourceTypeSetToStringSlice(intersection)
	err = d.Set("mapped_elements", mes)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(generateID(m.ID, mes))
	return
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	mID := d.Get("metaport_id").(string)
	rg, err := client.GetMetaport(ctx, c, mID)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing attachment of metaport %s because it's gone", mID)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	return attachmentToResource(d, rg)
}
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	mID := d.Get("metaport_id").(string)
	mes := client.ResourceTypeSetToStringSlice(d.Get("mapped_elements").(*schema.Set))
	m, err := client.AddMappedElementsToMetaport(ctx, c, mID, mes)
	if err != nil {
		return diag.FromErr(err)
	}
	return attachmentToResource(d, m)
}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	mID := d.Get("metaport_id").(string)
	mes := d.Get("mapped_elements").(*schema.Set)
	_, err := client.RemoveMappedElementsFromMetaport(ctx, c, mID, client.ResourceTypeSetToStringSlice(mes))
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	d.SetId("")
	return
}
