package access_control

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

const (
	description = `Device Access Control allows organizations to prevent clients from accessing unauthorized sites and services while they are connected to corporate resources. 
With this, organizations can block unfiltered/unmonitored client access when the users are connected to Proofpoint ZTNA.

When Device Access Control is enabled, outgoing traffic from the endpoint is allowed if it is either:
- Routed through Proofpoint ZTNA.
- Destined to an allowed service on the Internet.
Any other outgoing traffic is dropped.

~> **NOTE:** Proofpoint recommends the Device Access Control to be enabled only when a Web Security solution exists, and its IP ranges are allowed. Otherwise, users will fail to access the Internet when connected to Proofpoint ZTNA.
`
	applyToOrgDesc      = "Indicates whether this Access Control setting applies to the whole org. Note: This attribute overrides `apply_to_entities`."
	applyToEntitiesDesc = "Entities (users, groups or network elements) to be subjected to the Access Control."
	exemptEntitiesDesc  = "Entities (users, groups or network elements) which are exempt from the Access Control."
	allowedRoutesDesc   = "List of allowed IPv4 route CIDRs."
)

var excludedKeys = []string{"id"}

func accessControlRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	ac, err := client.GetAccessControl(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(ac, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ac.ID)
	return
}
func accessControlCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewAccessControl(d)
	ac, err := client.CreateAccessControl(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ac.ID)
	return accessControlRead(ctx, d, c)
}

func accessControlUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewAccessControl(d)
	ac, err := client.UpdateAccessControl(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ac.ID)
	return accessControlRead(ctx, d, c)
}

func accessControlDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteAccessControl(ctx, c, id)
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
