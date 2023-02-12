package client

import (
	"context"
	"encoding/json"
	"fmt"
	u "net/url"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	devicesEndpoint string = "v1/devices"
)

type Device struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description"`
	Enabled     *bool    `json:"enabled,omitempty"`
	OwnerID     string   `json:"owner_id,omitempty"`
	AutoAliases []string `json:"auto_aliases,omitempty"`
	Groups      []string `json:"groups,omitempty"`
	Tags        []Tag    `json:"tags,omitempty"`
	Aliases     []string `json:"aliases,omitempty"`
}

func NewDevice(d *schema.ResourceData) *Device {
	res := &Device{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}

	res.Description = d.Get("description").(string)
	enabled, exists := d.GetOkExists("enabled")
	if exists {
		enabled := enabled.(bool)
		res.Enabled = &enabled
	}
	if d.HasChange("owner_id") {
		res.OwnerID = d.Get("owner_id").(string)
	}
	return res
}

func parseDevice(resp []byte) (*Device, error) {
	device := &Device{}
	err := json.Unmarshal(resp, device)
	if err != nil {
		return nil, fmt.Errorf("could not parse device response: %v", err)
	}
	return device, nil
}

func CreateDevice(ctx context.Context, c *Client, dev *Device) (*Device, error) {
	devUrl := fmt.Sprintf("%s/%s", c.BaseURL, devicesEndpoint)
	body, err := json.Marshal(dev)
	if err != nil {
		return nil, fmt.Errorf("could not convert device to json: %v", err)
	}
	resp, err := c.Post(ctx, devUrl, body)
	if err != nil {
		return nil, err
	}
	return parseDevice(resp)
}

func UpdateDevice(ctx context.Context, c *Client, deviceId string, device *Device) (*Device, error) {
	deviceUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, devicesEndpoint, deviceId)
	body, err := json.Marshal(device)
	if err != nil {
		return nil, fmt.Errorf("could not convert device to json: %v", err)
	}
	resp, err := c.Patch(ctx, deviceUrl, body)
	if err != nil {
		return nil, err
	}
	return parseDevice(resp)
}

func GetDevice(ctx context.Context, c *Client, deviceID string) (*Device, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, devicesEndpoint, deviceID)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseDevice(resp)
}

func DeleteDevice(ctx context.Context, c *Client, deviceID string) (*Device, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, devicesEndpoint, deviceID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseDevice(resp)
}
