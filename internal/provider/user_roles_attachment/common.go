package user_roles_attachment

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

func generateID(uID string, gRoles []string) string {
	hash := 0
	for _, rID := range gRoles {
		hash += schema.HashString(rID)
	}
	return fmt.Sprintf("%s-%d", uID, hash)
}

func attachmentToResource(d *schema.ResourceData, uID string, r []string) (diags diag.Diagnostics) {
	err := d.Set("user_id", uID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("roles", r)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(generateID(uID, r))
	return
}

func readResource(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	uID := d.Get("user_id").(string)
	u, err := client.GetUserByID(ctx, c, uID)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			diags = append(diags, diag.Diagnostic{Severity: diag.Warning, Summary: fmt.Sprintf("user %s was not not found", uID)})
			return
		} else {
			return diag.FromErr(err)
		}
	}
	return attachmentToResource(d, uID, u.Roles)
}
func createResource(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	uID := d.Get("user_id").(string)
	r := client.ResourceTypeSetToStringSlice(d.Get("roles").(*schema.Set))
	roles, err := client.AssignRolesToUser(ctx, c, uID, r)
	if err != nil {
		return diag.FromErr(err)
	}
	return attachmentToResource(d, uID, roles)

}

func deleteResource(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	uID := d.Get("user_id").(string)
	_, err := client.AssignRolesToUser(ctx, c, uID, []string{})
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			diags = append(diags, diag.Diagnostic{Severity: diag.Warning, Summary: fmt.Sprintf("user %s was not not found", uID)})
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return
}
