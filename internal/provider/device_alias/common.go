package device_alias

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
	description = "DNS alias (FQDN) of the device."
)

func aliasToResource(d *schema.ResourceData, deviceID, alias string) diag.Diagnostics {
	var diags diag.Diagnostics
	err := d.Set("device_id", deviceID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = d.Set("alias", alias)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%s-%s", deviceID, alias))
	return diags
}

func deviceAliasRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	deviceID := d.Get("device_id").(string)
	alias := d.Get("alias").(string)
	exists, err := client.AliasExists(ctx, c, deviceID, alias)
	if err != nil {
		return diag.FromErr(err)
	}
	if !exists {
		log.Printf("Removing alias %s of device %s because it's gone", alias, deviceID)
		d.SetId("")
		return
	}
	return aliasToResource(d, deviceID, alias)
}
func deviceAliasCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	deviceID := d.Get("device_id").(string)
	alias := d.Get("alias").(string)
	err := client.AssignNetworkElementAlias(ctx, c, deviceID, alias)
	if err != nil {
		return diag.FromErr(err)
	}
	return aliasToResource(d, deviceID, alias)
}

func deviceAliasDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	deviceID := d.Get("device_id").(string)
	alias := d.Get("alias").(string)
	err := client.DeleteNetworkElementAlias(ctx, c, deviceID, alias)
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
