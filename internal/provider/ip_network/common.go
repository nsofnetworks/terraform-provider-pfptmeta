package ip_network

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

var excludedKeys = []string{"id", "org_id"}

const (
	description = "You can define ranges of IP addresses to be used as indicators of user location. " +
		"These ranges are intended for use as determining conditions in other resources."
	cirdsDesc     = "list of cidrs included in the network"
	countriesDesc = "list of countries included in the network"
)

func ipNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	cc, err := client.GetIPNetwork(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing ip network %s because it's gone", id)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId(cc.ID)
	err = client.MapResponseToResource(cc, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}
func ipNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewIPNetwork(d)
	cc, err := client.CreateIPNetwork(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cc.ID)
	err = client.MapResponseToResource(cc, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func ipNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewIPNetwork(d)
	cc, err := client.UpdateIPNetwork(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cc.ID)
	err = client.MapResponseToResource(cc, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func ipNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteIPNetwork(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return
}
