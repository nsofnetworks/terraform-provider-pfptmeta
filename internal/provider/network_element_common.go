package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

var networkElementExcludedKeys = []string{"id", "tags", "org_id"}

func setTags(neId string, d *schema.ResourceData, c *client.Client) error {
	rawTags := d.Get("tags").(map[string]interface{})
	tags := make([]*client.Tag, len(rawTags))
	index := 0
	for key, value := range rawTags {
		Tag := &client.Tag{
			Name:  key,
			Value: value.(string),
		}
		tags[index] = Tag
		index += 1
	}
	err := client.AssignNetworkElementTags(c, neId, tags)
	if err != nil {
		return err
	}
	return nil
}

func networkElementsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	networkElement, err := client.GetNetworkElement(c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	err = client.MapResponseToResource(networkElement, d, networkElementExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	tags := client.ConvertTagsListToMap(networkElement.Tags)
	err = d.Set("tags", tags)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(networkElement.ID)
	return diags
}
func networkElementCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	body := client.NewNetworkElementBody(d)
	networkElement, err := client.CreateNetworkElement(c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("tags") {
		err = setTags(networkElement.ID, d, c)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	networkElement, err = client.GetNetworkElement(c, networkElement.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(networkElement, d, networkElementExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	tags := client.ConvertTagsListToMap(networkElement.Tags)
	err = d.Set("tags", tags)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(networkElement.ID)
	return diags
}

func networkElementUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	networkElement, err := client.GetNetworkElement(c, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if networkElement == nil {
		d.SetId("")
	}
	if d.HasChanges("name", "description", "enabled", "mapped_subnets", "mapped_service") {
		body := client.NewNetworkElementBody(d)
		networkElement, err = client.UpdateNetworkElement(c, id, body)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("tags") {
		err = setTags(networkElement.ID, d, c)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	networkElement, err = client.GetNetworkElement(c, networkElement.ID)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(networkElement, d, networkElementExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	tags := client.ConvertTagsListToMap(networkElement.Tags)
	err = d.Set("tags", tags)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(networkElement.ID)
	return diags
}

func networkElementDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteNetworkElement(c, id)
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
