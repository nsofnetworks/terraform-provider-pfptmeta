package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"net/http"
)

const routingGroupsEndpoint string = "v1/routing_groups"

type RoutingGroup struct {
	ID                string   `json:"id,omitempty"`
	Name              string   `json:"name,omitempty"`
	Description       string   `json:"description"`
	MappedElementsIds []string `json:"mapped_elements_ids"`
	Sources           []string `json:"sources"`
	ExemptSources     []string `json:"exempt_sources"`
}

func NewRoutingGroup(d *schema.ResourceData) *RoutingGroup {
	res := &RoutingGroup{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)

	mes := d.Get("mapped_elements_ids")
	res.MappedElementsIds = ResourceTypeSetToStringSlice(mes.(*schema.Set))

	s := d.Get("sources")
	res.Sources = ResourceTypeSetToStringSlice(s.(*schema.Set))

	es := d.Get("exempt_sources")
	res.ExemptSources = ResourceTypeSetToStringSlice(es.(*schema.Set))
	return res
}

func parseRoutingGroup(resp *http.Response) (*RoutingGroup, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	pg := &RoutingGroup{}
	err = json.Unmarshal(body, pg)
	if err != nil {
		return nil, fmt.Errorf("could not parse routing group response: %v", err)
	}
	return pg, nil
}

func CreateRoutingGroup(ctx context.Context, c *Client, rg *RoutingGroup) (*RoutingGroup, error) {
	rgUrl := fmt.Sprintf("%s/%s", c.BaseURL, routingGroupsEndpoint)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert routing group to json: %v", err)
	}
	resp, err := c.Post(ctx, rgUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseRoutingGroup(resp)
}

func UpdateRoutingGroup(ctx context.Context, c *Client, rgID string, rg *RoutingGroup) (*RoutingGroup, error) {
	rgUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, routingGroupsEndpoint, rgID)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert routing group to json: %v", err)
	}
	resp, err := c.Patch(ctx, rgUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseRoutingGroup(resp)
}

func GetRoutingGroup(ctx context.Context, c *Client, rgID string) (*RoutingGroup, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, routingGroupsEndpoint, rgID)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseRoutingGroup(resp)
}

func DeleteRoutingGroup(ctx context.Context, c *Client, pgID string) (*RoutingGroup, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, routingGroupsEndpoint, pgID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseRoutingGroup(resp)
}

func AddMappedElementsToRoutingGroups(ctx context.Context, c *Client, rgID string, meIDs []string) (*RoutingGroup, error) {
	url := fmt.Sprintf("%s/%s/%s/add_mapped_elements", c.BaseURL, routingGroupsEndpoint, rgID)
	body := make(map[string][]string)
	body["mapped_element_ids"] = meIDs
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("could not convert mapped elements to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	return parseRoutingGroup(resp)
}

func RemoveMappedElementsFromRoutingGroups(ctx context.Context, c *Client, pgID string, meIDs []string) (*RoutingGroup, error) {
	url := fmt.Sprintf("%s/%s/%s/remove_mapped_elements", c.BaseURL, routingGroupsEndpoint, pgID)
	body := make(map[string][]string)
	body["mapped_element_ids"] = meIDs
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("could not convert mapped elements to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	return parseRoutingGroup(resp)
}
