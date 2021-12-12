package group_users_attachment

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

func generateID(gID string, users []string) string {
	hash := 0
	for _, uID := range users {
		hash += schema.HashString(uID)
	}
	return fmt.Sprintf("%s-%d", gID, hash)
}

func groupToUsersAttachmentResource(d *schema.ResourceData, g *client.Group) (diags diag.Diagnostics) {
	err := d.Set("group_id", g.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	gUsers := &schema.Set{F: schema.HashString}
	for _, i := range g.Users {
		gUsers.Add(i)
	}
	schemaUsers := d.Get("users").(*schema.Set)
	u := schema.NewSet(schema.HashString, schemaUsers.List())
	intersection := gUsers.Intersection(u)
	users := client.ResourceTypeSetToStringSlice(intersection)
	err = d.Set("users", users)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(generateID(g.ID, users))
	return
}

func readResource(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	gID := d.Get("group_id").(string)
	g, err := client.GetGroupById(c, gID)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			diags = append(diags, diag.Diagnostic{Severity: diag.Warning, Summary: fmt.Sprintf("group %s was not not found", gID)})
			return
		} else {
			return diag.FromErr(err)
		}
	}
	return groupToUsersAttachmentResource(d, g)
}
func createResource(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	gID := d.Get("group_id").(string)
	u := client.ResourceTypeSetToStringSlice(d.Get("users").(*schema.Set))
	err := client.AddUsersToGroup(c, gID, u)
	if err != nil {
		return diag.FromErr(err)
	}
	return readResource(nil, d, c)
}

func updateResource(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	gID := d.Get("group_id").(string)
	if d.HasChange("users") {
		before, after := d.GetChange("users")
		beforeSet, afterSet := before.(*schema.Set), after.(*schema.Set)
		toRemove := beforeSet.Difference(afterSet)
		toAdd := afterSet.Difference(beforeSet)
		if toRemove.Len() > 0 {
			err := client.RemoveUsersFromGroup(c, gID, client.ResourceTypeSetToStringSlice(toRemove))
			if err != nil {
				return append(diag.FromErr(err), readResource(nil, d, c)...)
			}
		}
		if toAdd.Len() > 0 {
			err := client.AddUsersToGroup(c, gID, client.ResourceTypeSetToStringSlice(toAdd))
			if err != nil {
				return append(diag.FromErr(err), readResource(nil, d, c)...)
			}
		}
	}
	return readResource(nil, d, c)
}

func deleteResource(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	gID := d.Get("group_id").(string)
	u := d.Get("users").(*schema.Set)
	err := client.RemoveUsersFromGroup(c, gID, client.ResourceTypeSetToStringSlice(u))
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
