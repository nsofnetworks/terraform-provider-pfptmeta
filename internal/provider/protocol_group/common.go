package protocol_group

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

var protocolGroupExcludedKeys = []string{"id"}

func protocolGroupCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	body := client.NewProtocolGroup(d)
	pg, err := client.CreateProtocolGroup(c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(pg.ID)
	err = client.MapResponseToResource(pg, d, protocolGroupExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func protocolGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	var pg *client.ProtocolGroup
	var err error
	if id, exists := d.GetOk("id"); exists {
		pg, err = client.GetProtocolGroupById(c, id.(string))
	}
	if name, exists := d.GetOk("name"); exists {
		pg, err = client.GetProtocolGroupByName(c, name.(string))
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
	if pg == nil {
		d.SetId("")
		return diags
	}
	err = client.MapResponseToResource(pg, d, protocolGroupExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(pg.ID)
	return diags
}

func protocolGroupUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	body := client.NewProtocolGroup(d)
	pg, err := client.UpdateProtocolGroup(c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(pg.ID)
	err = client.MapResponseToResource(pg, d, protocolGroupExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func protocolGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteProtocolGroup(c, id)
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
