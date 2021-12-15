package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"net/http"
	u "net/url"
)

const (
	networkElementsEndpoint string = "v1/network_elements"
)

type NetworkElementBody struct {
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description"`
	Enabled       bool     `json:"enabled"`
	MappedSubnets []string `json:"mapped_subnets,omitempty"`
	MappedService string   `json:"mapped_service,omitempty"`
	Platform      string   `json:"platform,omitempty"`
	OwnerID       string   `json:"owner_id,omitempty"`
}

func NewNetworkElementBody(d *schema.ResourceData) *NetworkElementBody {
	res := &NetworkElementBody{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}

	res.Description = d.Get("description").(string)

	res.Enabled = d.Get("enabled").(bool)

	if d.HasChange("mapped_subnets") {
		mappedSubnets := d.Get("mapped_subnets")
		listMappedSubnets := ResourceTypeSetToStringSlice(mappedSubnets.(*schema.Set))
		res.MappedSubnets = listMappedSubnets
	}
	if d.HasChange("mapped_service") {
		res.MappedService = d.Get("mapped_service").(string)
	}
	if d.HasChange("platform") {
		res.Platform = d.Get("platform").(string)
	}
	if d.HasChange("owner_id") {
		res.OwnerID = d.Get("owner_id").(string)
	}
	return res
}

type NetworkElementResponse struct {
	NetworkElementBody
	ID          string   `json:"id"`
	AutoAliases []string `json:"auto_aliases"`
	Groups      []string `json:"groups"`
	Type        string   `json:"type"`
	Tags        []Tag    `json:"tags,omitempty"`
	Aliases     []string `json:"aliases"`
}

func parseNetworkElement(resp *http.Response) (*NetworkElementResponse, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	networkElement := &NetworkElementResponse{}
	err = json.Unmarshal(body, networkElement)
	if err != nil {
		return nil, fmt.Errorf("could not parse network element response: %v", err)
	}
	return networkElement, nil
}

func CreateNetworkElement(ctx context.Context, c *Client, ne *NetworkElementBody) (*NetworkElementResponse, error) {
	neUrl := fmt.Sprintf("%s/%s", c.BaseURL, networkElementsEndpoint)
	body, err := json.Marshal(ne)
	if err != nil {
		return nil, fmt.Errorf("could not convert network element to json: %v", err)
	}
	resp, err := c.Post(ctx, neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseNetworkElement(resp)
}

func UpdateNetworkElement(ctx context.Context, c *Client, neId string, ne *NetworkElementBody) (*NetworkElementResponse, error) {
	neUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, networkElementsEndpoint, neId)
	body, err := json.Marshal(ne)
	if err != nil {
		return nil, fmt.Errorf("could not convert network element to json: %v", err)
	}
	resp, err := c.Patch(ctx, neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseNetworkElement(resp)
}

func GetNetworkElement(ctx context.Context, c *Client, neID string) (*NetworkElementResponse, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, networkElementsEndpoint, neID)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseNetworkElement(resp)
}

func DeleteNetworkElement(ctx context.Context, c *Client, neID string) (*NetworkElementResponse, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, networkElementsEndpoint, neID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseNetworkElement(resp)
}
