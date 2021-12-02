package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

var EnterpriseDNSExcludedKeys = []string{"id"}

func enterpriseDNSRead(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	ed, err := client.GetEnterpriseDNS(c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	err = client.MapResponseToResource(ed, d, EnterpriseDNSExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ed.ID)
	return
}
func enterpriseDNSCreate(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewEnterpriseDNS(d)
	ed, err := client.CreateEnterpriseDNS(c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(ed, d, EnterpriseDNSExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ed.ID)
	return
}

func enterpriseDNSUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewEnterpriseDNS(d)
	ed, err := client.UpdateEnterpriseDNS(c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(ed, d, EnterpriseDNSExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ed.ID)
	return
}

func enterpriseDNSDelete(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteEnterpriseDNS(c, id)
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
