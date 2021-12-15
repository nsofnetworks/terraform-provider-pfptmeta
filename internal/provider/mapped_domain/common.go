package mapped_domain

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

const (
	description = "DNS suffixes to be resolved within Mapped Subnet"
)

func mappedDomainToResource(d *schema.ResourceData, neID string, md *client.MappedDomain) (diags diag.Diagnostics) {
	err := client.MapResponseToResource(md, d, []string{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s-%s", neID, md.Name))
	return
}

func mappedDomainRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	neID := d.Get("network_element_id").(string)
	name := d.Get("name").(string)
	mappedDomain := d.Get("mapped_domain").(string)
	mdBody := &client.MappedDomain{MappedDomain: mappedDomain, Name: name}
	md, err := client.GetMappedDomain(ctx, c, neID, mdBody)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	return mappedDomainToResource(d, neID, md)
}
func mappedDomainCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	neID := d.Get("network_element_id").(string)
	name := d.Get("name").(string)
	mappedDomain := d.Get("mapped_domain").(string)
	mdBody := &client.MappedDomain{MappedDomain: mappedDomain, Name: name}
	md, err := client.SetMappedDomain(ctx, c, neID, mdBody)
	if err != nil {
		return diag.FromErr(err)
	}
	return mappedDomainToResource(d, neID, md)
}

func mappedDomainDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	neID := d.Get("network_element_id").(string)
	name := d.Get("name").(string)
	err := client.DeleteMappedDomain(ctx, c, neID, name)
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
