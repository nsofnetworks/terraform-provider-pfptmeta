package protocol_group

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

const (
	description   = "Protocol Groups are protocols and ports that must be included into granular policies."
	protocolsDesc = "A list of protocols"
	protoDesc     = "Protocol type, can be one of: tcp, udp, icmp"
)

var excludedKeys = []string{"id"}

func protocolGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	body := client.NewProtocolGroup(d)
	pg, err := client.CreateProtocolGroup(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(pg.ID)
	err = client.MapResponseToResource(pg, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func protocolGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	var pg *client.ProtocolGroup
	var err error
	if id, exists := d.GetOk("id"); exists {
		pg, err = client.GetProtocolGroupById(ctx, c, id.(string))
	}
	if name, exists := d.GetOk("name"); exists {
		pg, err = client.GetProtocolGroupByName(ctx, c, name.(string))
	}
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	if pg == nil {
		d.SetId("")
		return diags
	}
	err = client.MapResponseToResource(pg, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(pg.ID)
	return diags
}

func protocolGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	body := client.NewProtocolGroup(d)
	pg, err := client.UpdateProtocolGroup(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(pg.ID)
	err = client.MapResponseToResource(pg, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
func protocolGroupDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteProtocolGroup(ctx, c, id)
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
