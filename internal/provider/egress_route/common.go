package egress_route

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

const (
	description = `Egress routes allow traffic from a defined source (groups, users, devices or any other network element) to be routed to a specific MetaPort.
This can be useful for a variety of use cases including testing, security and compliance.
For example, you can define egress rules for SaaS applications that require access from a specific region for regulatory purposes.
You can set up egress rules to route traffic from a source to a MetaPort in your own data center or public cloud instance for service chaining or any type of traffic manipulation.`
	sourcesDesc       = "Entities (users, groups, devices or network elements) to be affected by the egress route (cannot be a Mapped Subnet if `via` is also a Mapped Subnet)."
	exemptSourcesDesc = "Entities (users, groups, devices or network elements) to be excluded from the egress route."
	destinationsDesc  = "Target hostnames or domains."
	viaDesc           = "Defines how the traffic will be routed:\n" +
		"	- **DIRECT**: Directs the traffic to egress from the same PoP it has entered. Use it to override other, less specific egress rules.\n" +
		"	- **Mapped Subnet ID**: Directs the traffic to egress via specified mapped subnet.\n" +
		"	- **Region**: Directs the traffic to egress from a specific region, see `location` data-source."
)

var excludedKeys = []string{"id"}

func egressRouteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	er, err := client.GetEgressRoute(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing egress_route %s because it's gone", id)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	err = client.MapResponseToResource(er, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(er.ID)
	return
}
func egressRouteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewEgressRoute(d)
	er, err := client.CreateEgressRoute(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(er, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(er.ID)
	return
}

func egressRouteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewEgressRoute(d)
	er, err := client.UpdateEgressRoute(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(er, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(er.ID)
	return
}

func egressRouteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteEgressRoute(ctx, c, id)
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
