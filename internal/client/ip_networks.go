package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	u "net/url"
)

const IPNetworksEndpoint = "v1/ip_networks"

type IPNetwork struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description"`
	Cidrs       []string `json:"cidrs"`
}

func NewIPNetwork(d *schema.ResourceData) *IPNetwork {
	res := &IPNetwork{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Cidrs = ConfigToStringSlice("cidrs", d)

	return res
}

func parseIPNetwork(resp []byte) (*IPNetwork, error) {
	in := &IPNetwork{}
	err := json.Unmarshal(resp, in)
	if err != nil {
		return nil, fmt.Errorf("could not parse ip network response: %v", err)
	}
	return in, nil
}

func CreateIPNetwork(ctx context.Context, c *Client, in *IPNetwork) (*IPNetwork, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, IPNetworksEndpoint)
	body, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("could not convert ip network to json: %v", err)
	}
	resp, err := c.Post(ctx, url, body)
	if err != nil {
		return nil, err
	}
	return parseIPNetwork(resp)
}

func UpdateIPNetwork(ctx context.Context, c *Client, inId string, in *IPNetwork) (*IPNetwork, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, IPNetworksEndpoint, inId)
	body, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("could not convert ip network to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, body)
	if err != nil {
		return nil, err
	}
	return parseIPNetwork(resp)
}

func GetIPNetwork(ctx context.Context, c *Client, inId string) (*IPNetwork, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, IPNetworksEndpoint, inId)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseIPNetwork(resp)
}

func DeleteIPNetwork(ctx context.Context, c *Client, inId string) (*IPNetwork, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, IPNetworksEndpoint, inId)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseIPNetwork(resp)
}
