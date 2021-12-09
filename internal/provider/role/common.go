package role

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

var roleExcludedKeys = []string{"id"}

func roleCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	body := client.NewRole(d)
	r, err := client.CreateRole(c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(r.ID)
	err = client.MapResponseToResource(r, d, roleExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func roleRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	var r *client.Role
	var err error
	if id, exists := d.GetOk("id"); exists {
		r, err = client.GetRoleByID(c, id.(string))
	} else {
		if name, exists := d.GetOk("name"); exists {
			r, err = client.GetRoleByName(c, name.(string))
		}
	}
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	if r == nil {
		d.SetId("")
		return diags
	}
	err = client.MapResponseToResource(r, d, roleExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(r.ID)
	return diags
}

func roleUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	body := client.NewRole(d)
	r, err := client.UpdateRole(c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(r.ID)
	err = client.MapResponseToResource(r, d, roleExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func roleDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteRole(c, id)
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
