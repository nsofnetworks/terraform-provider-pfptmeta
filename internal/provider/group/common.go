package group

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

const (
	description    = "Groups represent a collection of users, typically belong to a common department or share same privileges in the organization."
	expressionDesc = "Allows grouping entities by their tags. Filtering by tag value is also supported if provided. " +
		"Supported operations: AND, OR, NOT, parenthesis."
)

var excludedKeys = []string{"id", "expression", "roles", "users"}

func groupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewGroup(d)
	g, err := client.CreateGroup(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return groupToResource(d, g)
}

func groupToResource(d *schema.ResourceData, g *client.Group) (diags diag.Diagnostics) {
	d.SetId(g.ID)
	err := client.MapResponseToResource(g, d, excludedKeys)
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

func groupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	var g *client.Group
	var err error
	id, exists := d.GetOk("id")
	if exists {
		g, err = client.GetGroupById(ctx, c, id.(string))
	}
	if name, exists := d.GetOk("name"); exists {
		g, err = client.GetGroupByName(ctx, c, name.(string))
		if g == nil {
			log.Printf("[WARN] Removing group %s because it's gone", name)
			d.SetId("")
			return diags
		}
	}
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing group %s because it's gone", id)
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	err = client.MapResponseToResource(g, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(g.ID)
	return diags
}

func groupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	id := d.Id()
	body := client.NewGroup(d)
	g, err := client.UpdateGroup(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(g.ID)
	return groupToResource(d, g)

}
func groupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteGroup(ctx, c, id)
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
