package network_element

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

const (
	description = "Network elements comprise devices, mapped subnets and mapped services. \n" +
		"- Creating this resource with `mapped_subnets` generates a Mapped Subnet-type network element..\n" +
		"- Creating this resource with `mapped_service` generates a Mapped Service-type network element.\n" +
		"- Creating this resource with `owner_id` and `platform` generates a Device-type network element.\n"
	tagsDesc          = "Key/value attributes for combining elements together into Smart Groups, and placed as targets or sources in Policies"
	mappedSubnetsDesc = "CIDRs that will be mapped to the subnet"
	enabledDesc       = "Not allowed for mapped service and mapped domain"
)

var excludedKeys = []string{"id", "tags", "aliases"}

func networkElementsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	networkElement, err := client.GetNetworkElement(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(networkElement, d, excludedKeys)
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
func networkElementCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewNetworkElementBody(d)
	networkElement, err := client.CreateNetworkElement(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(networkElement.ID)
	err = client.MapResponseToResource(networkElement, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return updateTags(ctx, d, networkElement, c)
}

func networkElementUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	id := d.Id()
	networkElement, err := client.GetNetworkElement(ctx, c, id)
	if err != nil {
		return diag.FromErr(err)
	}
	if networkElement == nil {
		d.SetId("")
	}
	if d.HasChanges("name", "description", "enabled", "mapped_subnets", "mapped_service") {
		body := client.NewNetworkElementBody(d)
		networkElement, err = client.UpdateNetworkElement(ctx, c, id, body)
		if err != nil {
			return diag.FromErr(err)
		}
		err = client.MapResponseToResource(networkElement, d, excludedKeys)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return updateTags(ctx, d, networkElement, c)
}

func networkElementDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteNetworkElement(ctx, c, id)
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

func updateTags(ctx context.Context, d *schema.ResourceData, ne *client.NetworkElementResponse, c *client.Client) (diags diag.Diagnostics) {
	if d.HasChange("tags") {
		tags := client.NewTags(d)
		err := client.AssignTagsToResource(ctx, c, ne.ID, "network_elements", tags)
		if err != nil {
			return diag.FromErr(err)
		}
		networkElementsRead(ctx, d, c)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return
}
