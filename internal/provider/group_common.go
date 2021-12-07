package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

var GroupExcludedKeys = []string{"id", "expression"}

func groupCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewGroup(d)
	g, err := client.CreateGroup(c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return groupToResource(d, g)
}

func groupToResource(d *schema.ResourceData, g *client.Group) (diags diag.Diagnostics) {
	d.SetId(g.ID)
	err := client.MapResponseToResource(g, d, GroupExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	if g.Expression == nil {
		d.Set("expression", "")
	} else {
		err = d.Set("expression", *g.Expression)
	}
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func groupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	var pg *client.Group
	var err error
	if id, exists := d.GetOk("id"); exists {
		pg, err = client.GetGroupById(c, id.(string))
	}
	if name, exists := d.GetOk("name"); exists {
		pg, err = client.GetGroupByName(c, name.(string))
	} else if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	if pg == nil {
		d.SetId("")
		return diags
	}
	err = client.MapResponseToResource(pg, d, GroupExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(pg.ID)
	return diags
}

func groupUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	id := d.Id()
	body := client.NewGroup(d)
	g, err := client.UpdateGroup(c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(g.ID)
	return groupToResource(d, g)

}
func groupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteGroup(c, id)
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
