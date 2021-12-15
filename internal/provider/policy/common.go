package policy

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

const (
	description = "Policies bind network elements (devices, services and subnets) to users, " +
		"groups and to other network elements and define access direction and connections."
	destinationsDesc   = "Entities (users, groups or network elements) to which the access is granted to."
	sourcesDesc        = "Entities (users, groups or network elements) to be authorized to access the application defined in this policy."
	exemptSourcesDesc  = "Entities (users, groups or network elements) to be excluded from accessing the application defined in this policy."
	protocolGroupsDesc = "Protocol groups that restrict the protocols or TCP/UDP ports for this policy"
)

var policyExcludedKeys = []string{"id"}

func policyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	body := client.NewPolicy(d)
	r, err := client.CreatePolicy(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(r.ID)
	err = client.MapResponseToResource(r, d, policyExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func policyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Get("id").(string)
	rg, err := client.GetPolicy(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	err = client.MapResponseToResource(rg, d, policyExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(rg.ID)
	return diags
}

func policyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	body := client.NewPolicy(d)
	r, err := client.UpdatePolicy(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(r.ID)
	err = client.MapResponseToResource(r, d, policyExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func policyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeletePolicy(ctx, c, id)
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
