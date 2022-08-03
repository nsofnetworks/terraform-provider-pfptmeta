package network_element_alias

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
)

const (
	description = "DNS alias (FQDN) of the network element. valid for network element of type Device, Native Service and Mapped Service."
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

func networkElementsAliasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	neID := d.Get("network_element_id").(string)
	alias := d.Get("alias").(string)
	exists, err := client.AliasExists(ctx, c, neID, alias)
	if err != nil {
		return diag.FromErr(err)
	}
	if !exists {
		log.Printf("Removing alias %s of network element %s because it's gone", alias, neID)
		d.SetId("")
		return
	}
	return aliasToResource(d, neID, alias)
}
func networkElementAliasCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	neID := d.Get("network_element_id").(string)
	alias := d.Get("alias").(string)
	err := client.AssignAlias(ctx, c, neID, alias)
	if err != nil {
		return diag.FromErr(err)
	}
	return aliasToResource(d, neID, alias)
}

func networkElementAliasDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	neID := d.Get("network_element_id").(string)
	alias := d.Get("alias").(string)
	err := client.DeleteAlias(ctx, c, neID, alias)
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
