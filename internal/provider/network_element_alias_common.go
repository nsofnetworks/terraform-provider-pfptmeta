package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

func aliasToResource(d *schema.ResourceData, neID, alias string) diag.Diagnostics {
	var diags diag.Diagnostics
	err := d.Set("network_element_id", neID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("alias", alias)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s-%s", neID, alias))
	return diags
}

func networkElementsAliasRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	neID := d.Get("network_element_id").(string)
	alias := d.Get("alias").(string)
	exists, err := client.AliasExists(c, neID, alias)
	if err != nil {
		return diag.FromErr(err)
	}
	if !exists {
		d.SetId("")
	}
	return aliasToResource(d, neID, alias)
}
func networkElementAliasCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	neID := d.Get("network_element_id").(string)
	alias := d.Get("alias").(string)
	err := client.AssignNetworkElementAlias(c, neID, alias)
	if err != nil {
		return diag.FromErr(err)
	}
	return aliasToResource(d, neID, alias)
}

func networkElementAliasDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	neID := d.Get("network_element_id").(string)
	alias := d.Get("alias").(string)
	err := client.DeleteNetworkElementAlias(c, neID, alias)
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
