package tunnel

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
)

var excludedKeys = []string{"id", "gre_config"}

const (
	description = "Tunnels represent the origin of the connection from " +
		"the customer’s site to the Proofpoint cloud."
)

func greTunnelConfigFromResource(d *schema.ResourceData) *client.GreTunnelConfig {
	g := d.Get("gre_config").([]interface{})
	if len(g) == 0 {
		return nil
	}
	o := g[0].(map[string]interface{})
	sips := o["source_ips"]
	ips := client.ResourceTypeSetToStringSlice(sips.(*schema.Set))
	return &client.GreTunnelConfig{SourceIps: ips}
}

func tunnelFromResource(d *schema.ResourceData) *client.Tunnel {
	res := &client.Tunnel{}
	if d.HasChange("name") {
		name := d.Get("name")
		res.Name = name.(string)
	}

	res.Description = d.Get("description").(string)
	res.GreConfig = greTunnelConfigFromResource(d)

	enabled := d.Get("enabled").(bool)
	res.Enabled = &enabled

	return res
}

func tunnelToResource(d *schema.ResourceData, t *client.Tunnel) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId(t.ID)
	if err := client.MapResponseToResource(t, d, excludedKeys); err != nil {
		return diag.FromErr(err)
	}
	if t.GreConfig != nil {
		g := []map[string]interface{}{
			{
				"source_ips": t.GreConfig.SourceIps,
			},
		}
		if err := d.Set("gre_config", g); err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

func tunnelRead(keepInStateOnMissing bool, ctx context.Context,
	d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	var err error
	var t *client.Tunnel
	id, exists := d.GetOk("id")
	if exists {
		t, err = client.GetTunnel(ctx, c, id.(string))
	}
	name, exists := d.GetOk("name")
	if exists && t == nil {
		t, err = client.GetTunnelByName(ctx, c, name.(string))
	}
	if err != nil {
		if keepInStateOnMissing {
			return diag.FromErr(err)
		}
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing tunnel %s because it's gone", id)
			d.SetId("")
			return diags
		} else if strings.HasPrefix(err.Error(), "could not find tunnel") {
			log.Printf("[WARN] Removing tunnel %s because it's gone", name)
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	return tunnelToResource(d, t)
}

func getTunnelRead(keepInStateOnMissing bool) schema.ReadContextFunc {
	return func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		return tunnelRead(keepInStateOnMissing, ctx, d, meta)
	}
}

func tunnelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	var t *client.Tunnel

	tTmpl := tunnelFromResource(d)
	t, err := client.CreateTunnel(ctx, c, tTmpl)
	if err != nil {
		return diag.FromErr(err)
	}
	if tTmpl.GreConfig != nil {
		t, err = client.AddGreSourceIpsToTunnel(ctx, c, t.ID,
			tTmpl.GreConfig.SourceIps)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return tunnelToResource(d, t)
}

func tunnelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	var t *client.Tunnel

	id := d.Id()
	tTmpl := tunnelFromResource(d)
	t, err := client.UpdateTunnel(ctx, c, id, tTmpl)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("gre_config") {
		before, after := d.GetChange("gre_config.0.source_ips")
		beforeSet, afterSet := before.(*schema.Set), after.(*schema.Set)
		toRemove := beforeSet.Difference(afterSet)
		toAdd := afterSet.Difference(beforeSet)
		if toRemove.Len() > 0 {
			t, err = client.RemoveGreSourceIpsFromTunnel(ctx, c, id,
				client.ResourceTypeSetToStringSlice(toRemove))
			if err != nil {
				return diag.FromErr(err)
			}
		}
		if toAdd.Len() > 0 {
			t, err = client.AddGreSourceIpsToTunnel(ctx, c, id,
				client.ResourceTypeSetToStringSlice(toAdd))
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	return tunnelToResource(d, t)
}

func tunnelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteTunnel(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	return diags
}
