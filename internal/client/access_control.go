package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	u "net/url"
)

const accessControlEndpoint = "v1/access_controls"

type AccessControl struct {
	ID              string   `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description"`
	Enabled         bool     `json:"enabled"`
	ApplyToOrg      bool     `json:"apply_to_org"`
	ApplyToEntities []string `json:"apply_to_entities"`
	ExemptEntities  []string `json:"exempt_entities"`
	AllowedRoutes   []string `json:"allowed_routes"`
}

func NewAccessControl(d *schema.ResourceData) *AccessControl {
	res := &AccessControl{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Enabled = d.Get("enabled").(bool)
	res.ApplyToOrg = d.Get("apply_to_org").(bool)
	res.ApplyToEntities = ConfigToStringSlice("apply_to_entities", d)
	res.ExemptEntities = ConfigToStringSlice("exempt_entities", d)
	res.AllowedRoutes = ConfigToStringSlice("allowed_routes", d)
	return res
}

func parseAccessControl(resp []byte) (*AccessControl, error) {
	e := &AccessControl{}
	err := json.Unmarshal(resp, e)
	if err != nil {
		return nil, fmt.Errorf("could not parse access control response: %v", err)
	}
	return e, nil
}

func CreateAccessControl(ctx context.Context, c *Client, e *AccessControl) (*AccessControl, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, accessControlEndpoint)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("could not convert access control to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseAccessControl(resp)
}

func GetAccessControl(ctx context.Context, c *Client, eID string) (*AccessControl, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, accessControlEndpoint, eID)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseAccessControl(resp)
}

func UpdateAccessControl(ctx context.Context, c *Client, eID string, e *AccessControl) (*AccessControl, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, accessControlEndpoint, eID)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("could not convert access control to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseAccessControl(resp)
}

func DeleteAccessControl(ctx context.Context, c *Client, mID string) (*AccessControl, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, accessControlEndpoint, mID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseAccessControl(resp)
}
