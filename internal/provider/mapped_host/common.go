package mapped_host

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

func mappedHostToResource(d *schema.ResourceData, neID string, mh *client.MappedHost) (diags diag.Diagnostics) {
	err := client.MapResponseToResource(mh, d, []string{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s-%s", neID, mh.Name))
	return
}

func mappedHostRead(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	neID := d.Get("network_element_id").(string)
	name := d.Get("name").(string)
	mappedHost := d.Get("mapped_host").(string)
	mhBody := &client.MappedHost{MappedHost: mappedHost, Name: name}
	mh, err := client.GetMappedHost(c, neID, mhBody)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	return mappedHostToResource(d, neID, mh)
}
func mappedHostCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	neID := d.Get("network_element_id").(string)
	name := d.Get("name").(string)
	mappedHost := d.Get("mapped_host").(string)
	mhBody := &client.MappedHost{MappedHost: mappedHost, Name: name}
	mh, err := client.SetMappedHost(c, neID, mhBody)
	if err != nil {
		return diag.FromErr(err)
	}
	return mappedHostToResource(d, neID, mh)
}

func mappedHostDelete(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	neID := d.Get("network_element_id").(string)
	name := d.Get("name").(string)
	err := client.DeleteMappedHost(c, neID, name)
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
