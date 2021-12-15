package group_roles_attachment

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

func generateID(gID string, gRoles []string) string {
	hash := 0
	for _, rID := range gRoles {
		hash += schema.HashString(rID)
	}
	return fmt.Sprintf("%s-%d", gID, hash)
}

func attachmentToResource(d *schema.ResourceData, gID string, r []string) (diags diag.Diagnostics) {
	err := d.Set("group_id", gID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("roles", r)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(generateID(gID, r))
	return
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	gID := d.Get("group_id").(string)
	g, err := client.GetGroupById(ctx, c, gID)
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
	return attachmentToResource(d, gID, g.Roles)
}
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	gID := d.Get("group_id").(string)
	r := client.ResourceTypeSetToStringSlice(d.Get("roles").(*schema.Set))
	roles, err := client.AssignRolesToGroup(ctx, c, gID, r)
	if err != nil {
		return diag.FromErr(err)
	}
	return attachmentToResource(d, gID, roles)

}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	gID := d.Get("group_id").(string)
	_, err := client.AssignRolesToGroup(ctx, c, gID, []string{})
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			diags = append(diags, diag.Diagnostic{Severity: diag.Warning, Summary: fmt.Sprintf("group %s was not not found", gID)})
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return
}
