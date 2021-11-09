package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
	"sync"
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
	c := meta.(*client.Client)

	body := client.NewNetworkElementBody(d)
	networkElement, err := client.CreateNetworkElement(c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(networkElement.ID)
	err = client.MapResponseToResource(networkElement, d, networkElementExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return updateExpandedAttributes(d, networkElement, c)
}

func networkElementUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		err = client.MapResponseToResource(networkElement, d, networkElementExcludedKeys)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return updateExpandedAttributes(d, networkElement, c)
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

func updateMappedDomains(neID string, d *schema.ResourceData, c *client.Client) diag.Diagnostics {
	var diags diag.Diagnostics
	oldMd, newMd := d.GetChange("mapped_domains")
	oldMdSet := oldMd.(*schema.Set)
	newMdSet := newMd.(*schema.Set)
	toDelete := oldMdSet.Difference(newMdSet)
	toWrite := newMdSet.Difference(oldMdSet)
	mdsToDelete := parseMappedDomains(toDelete)
	mdsToWrite := parseMappedDomains(toWrite)
	var wg sync.WaitGroup
	wg.Add(toDelete.Len() + toWrite.Len())
	diagsChan := make(chan diag.Diagnostics, toDelete.Len()+toWrite.Len())
	for _, md := range mdsToDelete {
		md := md
		go func() {
			defer wg.Done()
			var diags diag.Diagnostics
			err := client.DeleteMappedDomain(c, neID, md.Name)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
			diagsChan <- diags
		}()
	}
	for _, md := range mdsToWrite {
		md := md
		go func() {
			defer wg.Done()
			var diags diag.Diagnostics
			err := client.SetMappedDomain(c, neID, md)
			if err != nil {
				diags = append(diags, diag.FromErr(err)...)
			}
			diagsChan <- diags
		}()
	}
	wg.Wait()
	close(diagsChan)
	diags = append(diags, <-diagsChan...)
	diags = append(diags, networkElementsRead(nil, d, c)...)
	return diags
}

func parseMappedDomains(mds *schema.Set) []*client.MappedDomain {
	if mds.Len() == 0 {
		return nil
	}
	resp := make([]*client.MappedDomain, mds.Len())
	for i, v := range mds.List() {
		md := v.(map[string]interface{})
		resp[i] = &client.MappedDomain{Name: md["name"].(string), MappedDomain: md["mapped_domain"].(string)}
	}
	return resp
}

func updateExpandedAttributes(d *schema.ResourceData, networkElement *client.NetworkElementResponse, c *client.Client) diag.Diagnostics {
	var diags diag.Diagnostics
	if d.HasChange("tags") {
		err := setTags(networkElement.ID, d, c)
		networkElementsRead(nil, d, c)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("mapped_domains") {
		diags = append(diags, updateMappedDomains(networkElement.ID, d, c)...)
	}
	return diags
}