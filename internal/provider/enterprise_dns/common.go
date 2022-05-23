package enterprise_dns

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

const (
	description = "Enterprise DNS provides integration with global, enterprise DNS servers, " +
		"allowing resolution of FQDNs for domains that are in different locations/datacenters."
	mappedDomainsDesc = "DNS suffixes to be resolved within the enterprise DNS server"
	mappedDomainDesc  = "Proofpoint DNS Suffix"
	mdNameDescription = "Enterprise DNS server DNS suffix"
)

var excludedKeys = []string{"id"}

func enterpriseDNSRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	ed, err := client.GetEnterpriseDNS(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing enterprise_dns %s because it's gone", id)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	err = client.MapResponseToResource(ed, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ed.ID)
	return
}
func enterpriseDNSCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewEnterpriseDNS(d)
	ed, err := client.CreateEnterpriseDNS(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(ed, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ed.ID)
	return
}

func enterpriseDNSUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewEnterpriseDNS(d)
	ed, err := client.UpdateEnterpriseDNS(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(ed, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ed.ID)
	return
}

func enterpriseDNSDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteEnterpriseDNS(ctx, c, id)
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
