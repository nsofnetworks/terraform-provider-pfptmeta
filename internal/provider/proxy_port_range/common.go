package proxy_port_range

import (
	"context"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
)

var excludedKeys = []string{"id", "org_id"}

const (
	description = "Administrators can define communication ports for Web Security traffic over HTTP/S. The following port range is supported: 1 – 65535."
	proto       = "Protocol for this proxy port range. ENUM: `HTTPS`,`HTTP`"
	from_port   = "Start port for this proxy port range"
	to_port     = "End port for this proxy port range"
)

func proxyPortRangeToResource(d *schema.ResourceData, ppr *client.ProxyPortRange) (diags diag.Diagnostics) {
	d.SetId(ppr.ID)
	err := client.MapResponseToResource(ppr, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func proxyPortRangeRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	ppr, err := client.GetProxyPortRange(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing proxy port range %s because it's gone", id)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	return proxyPortRangeToResource(d, ppr)
}
func proxyPortRangeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewProxyPortRange(d)
	ppr, err := client.CreateProxyPortRange(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return proxyPortRangeToResource(d, ppr)
}

func proxyPortRangeUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewProxyPortRange(d)
	ppr, err := client.UpdateProxyPortRange(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return proxyPortRangeToResource(d, ppr)
}

func proxyPortRangeDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteProxyPortRange(ctx, c, id)
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
