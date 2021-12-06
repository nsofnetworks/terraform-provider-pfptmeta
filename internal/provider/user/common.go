package user

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

const (
	description = `Users are individuals in the organization, which their access and privileges to the organization resources are determined based on the security access policies.

Users can own multiple devices, each with a dedicated certificate that is used to identify the user’s device in the system.

Users can be provisioned in the system either by locally creating the users and groups from the Admin portal or API. Another, more common practice, is to provision users from an organization directory service, by way of SCIM or LDAP protocols.`
	givenNameDesc  = "User first name"
	familyNameDesc = "User family name"
	tagsDesc       = "Key/value attributes to be used for combining elements together into Smart Groups, and placed as targets or sources in Policies"
)

var userExcludedKeys = []string{"id", "tags"}

func userToResource(u *client.User, d *schema.ResourceData) (diags diag.Diagnostics) {
	err := client.MapResponseToResource(u, d, userExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	tags := client.ConvertTagsListToMap(u.Tags)
	err = d.Set("tags", tags)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(u.ID)
	return
}

func updateUserTags(d *schema.ResourceData, u *client.User, c *client.Client) (diags diag.Diagnostics) {
	if d.HasChange("tags") {
		tags := client.NewTags(d)
		err := client.AssignTagsToResource(c, u.ID, "users", tags)
		if err != nil {
			return diag.FromErr(err)
		}
		userRead(nil, d, c)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return
}

func userRead(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	var u *client.User
	var err error
	if id, exists := d.GetOk("id"); exists {
		u, err = client.GetUserByID(c, id.(string))
	} else {
		if email, exists := d.GetOk("email"); exists {
			u, err = client.GetUserByEmail(c, email.(string))
		}
	}
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	if u == nil {
		d.SetId("")
		return
	}
	return userToResource(u, d)
}

func userCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewUser(d)
	u, err := client.CreateUser(c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(u.ID)
	return updateUserTags(d, u, c)
}

func userUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewUser(d)
	u, err := client.UpdateUser(c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return updateUserTags(d, u, c)
}

func userDelete(_ context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteUser(c, id)
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
