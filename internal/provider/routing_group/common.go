package routing_group

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

const (
	description = "User routing groups is a logical group of subnets which co-exist at the same time. " +
		"These are independent subnets, which can be the same or contain overlapping IP addresses that can be used without conflicting with each other."
	mappedElementIdsDesc = "Mapped subnets and services that belong to this routing group."
	sourcesDesc          = "Users, groups or devices whose traffic will be routed."
	exemptSourcesDesc    = "Users, groups or devices whose traffic will not be routed."
)

var excludedKeys = []string{"id"}

func routingGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	body := client.NewRoutingGroup(d)
	r, err := client.CreateRoutingGroup(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(r.ID)
	err = client.MapResponseToResource(r, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func routingGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Get("id").(string)
	rg, err := client.GetRoutingGroup(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing routing group %s because it's gone", id)
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	err = client.MapResponseToResource(rg, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(rg.ID)
	return diags
}

func routingGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	body := client.NewRoutingGroup(d)
	r, err := client.UpdateRoutingGroup(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(r.ID)
	err = client.MapResponseToResource(r, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func routingGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteRoutingGroup(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}
