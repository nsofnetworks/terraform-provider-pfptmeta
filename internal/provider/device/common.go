package device

import (
	"context"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
)

const (
	description = "When a user is onboarded to the Proofpoint NaaS platform via Proofpoint Agent,\n" +
		"the user identity is bound to the device of the logging request, and a certificate is issued to this machine."
	tagsDesc = "Key/value attributes for combining elements together into Smart Groups, and placed as targets or sources in Policies"
)

var excludedKeys = []string{"id", "tags", "aliases"}

func deviceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	device, err := client.GetDevice(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing device %s because it's gone", id)
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	err = client.MapResponseToResource(device, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	handleEmptyListAttributes(d, device)
	tags := client.ConvertTagsListToMap(device.Tags)
	err = d.Set("tags", tags)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(device.ID)
	return diags
}

func deviceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewDevice(d)
	device, err := client.CreateDevice(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(device.ID)
	err = client.MapResponseToResource(device, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	handleEmptyListAttributes(d, device)
	return updateTags(ctx, d, device, c)
}

func deviceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	id := d.Id()
	device, err := client.GetDevice(ctx, c, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if device == nil {
		d.SetId("")
	}
	if d.HasChanges("name", "description", "enabled") {
		body := client.NewDevice(d)
		device, err = client.UpdateDevice(ctx, c, id, body)
		if err != nil {
			return diag.FromErr(err)
		}
		err = client.MapResponseToResource(device, d, excludedKeys)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	handleEmptyListAttributes(d, device)
	return updateTags(ctx, d, device, c)
}

func deviceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteDevice(ctx, c, id)
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

func updateTags(ctx context.Context, d *schema.ResourceData, device *client.Device, c *client.Client) (diags diag.Diagnostics) {
	if d.HasChange("tags") {
		tags := client.NewTags(d)
		err := client.AssignTagsToResource(ctx, c, device.ID, "devices", tags)
		if err != nil {
			return diag.FromErr(err)
		}
		deviceRead(ctx, d, c)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return
}

func handleEmptyListAttributes(d *schema.ResourceData, device *client.Device) {
	if device.AutoAliases == nil || len(device.AutoAliases) == 0 {
		d.Set("auto_aliases", make([]string, 0))
	}
	if device.Groups == nil || len(device.Groups) == 0 {
		d.Set("groups", make([]string, 0))
	}
}
