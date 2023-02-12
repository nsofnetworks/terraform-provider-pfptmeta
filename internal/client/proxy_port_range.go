package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ProxyPortRangeEndpoint = "v1/proxy_port_ranges"

type ProxyPortRange struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Proto       string `json:"proto"`
	FromPort    int    `json:"from_port"`
	ToPort      int    `json:"to_port"`
	ReadOnly    bool   `json:"read_only,omitempty"`
}

func NewProxyPortRange(d *schema.ResourceData) *ProxyPortRange {
	res := &ProxyPortRange{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Proto = d.Get("proto").(string)
	res.FromPort = d.Get("from_port").(int)
	res.ToPort = d.Get("to_port").(int)
	return res
}

func parseProxyPortRange(resp []byte) (*ProxyPortRange, error) {
	ppr := &ProxyPortRange{}
	err := json.Unmarshal(resp, ppr)
	if err != nil {
		return nil, fmt.Errorf("could not parse proxy port range response: %v", err)
	}
	return ppr, nil
}

func CreateProxyPortRange(ctx context.Context, c *Client, ppr *ProxyPortRange) (*ProxyPortRange, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, ProxyPortRangeEndpoint)
	body, err := json.Marshal(ppr)
	if err != nil {
		return nil, fmt.Errorf("could not convert proxy port to json: %v", err)
	}
	resp, err := c.Post(ctx, url, body)
	if err != nil {
		return nil, err
	}
	return parseProxyPortRange(resp)
}

func UpdateProxyPortRange(ctx context.Context, c *Client, pprId string, ppr *ProxyPortRange) (*ProxyPortRange, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, ProxyPortRangeEndpoint, pprId)
	body, err := json.Marshal(ppr)
	if err != nil {
		return nil, fmt.Errorf("could not convert proxy port to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, body)
	if err != nil {
		return nil, err
	}
	return parseProxyPortRange(resp)
}

func GetProxyPortRange(ctx context.Context, c *Client, pprId string) (*ProxyPortRange, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, ProxyPortRangeEndpoint, pprId)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseProxyPortRange(resp)
}

func DeleteProxyPortRange(ctx context.Context, c *Client, pprId string) (*ProxyPortRange, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, ProxyPortRangeEndpoint, pprId)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseProxyPortRange(resp)
}
