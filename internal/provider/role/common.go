package role

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

const (
	description           = "Roles define operations on the enterprise network, such as adding and removing users, defining security policies, etc."
	applyToOrgsDesc       = "indicates which orgs this role applies to. By default, it is applied to current org."
	privilegesDesc        = "Privileges to be assigned to the new role. It has the following structure - `resource:read/write` For example, metaports:read etc."
	subOrgsExpressionDesc = "Allows grouping of entities according to their tags. Filtering by tag value is also supported, if provided. Supported operations: AND, OR, XOR, parenthesis."
)

var excludedKeys = []string{"id", "privileges"}

func roleToResource(r *client.Role, d *schema.ResourceData) (diags diag.Diagnostics) {
	d.SetId(r.ID)
	err := client.MapResponseToResource(r, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	origPrivs := schema.NewSet(schema.HashString, d.Get("privileges").(*schema.Set).List())
	newPrivs := &schema.Set{F: schema.HashString}
	for _, i := range r.Privileges {
		newPrivs.Add(i)
	}
	err = d.Set("privileges", origPrivs.Intersection(newPrivs))
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func roleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewRole(d)
	r, err := client.CreateRole(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return roleToResource(r, d)
}
func roleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	var r *client.Role
	var err error
	if id, exists := d.GetOk("id"); exists {
		r, err = client.GetRoleByID(ctx, c, id.(string))
	} else {
		if name, exists := d.GetOk("name"); exists {
			r, err = client.GetRoleByName(ctx, c, name.(string))
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
	return roleToResource(r, d)
}

func roleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewRole(d)
	r, err := client.UpdateRole(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return roleToResource(r, d)
}

func roleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteRole(ctx, c, id)
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
