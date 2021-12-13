package metaport

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

var excludedKeys = []string{"id"}

func metaportRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	m, err := client.GetMetaport(c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	err = client.MapResponseToResource(m, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(m.ID)
	return diags
}

func metaportCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	body := client.NewMetaport(d)
	m, err := client.CreateMetaport(c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(m.ID)
	err = client.MapResponseToResource(m, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func metaportUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewMetaport(d)
	m, err := client.UpdateMetaport(c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(m, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func metaportDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteMetaport(c, id)
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
