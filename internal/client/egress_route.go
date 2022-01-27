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

const egressRouteEndpoint = "v1/egress_routes"

type EgressRoute struct {
	ID            string   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description"`
	Destinations  []string `json:"destinations"`
	Enabled       bool     `json:"enabled"`
	ExemptSources []string `json:"exempt_sources"`
	Sources       []string `json:"sources"`
	Via           string   `json:"via"`
}

func NewEgressRoute(d *schema.ResourceData) *EgressRoute {
	res := &EgressRoute{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)

	res.Enabled = d.Get("enabled").(bool)

	res.Destinations = ConfigToStringSlice("destinations", d)
	res.ExemptSources = ConfigToStringSlice("exempt_sources", d)
	res.Sources = ConfigToStringSlice("sources", d)
	res.Via = d.Get("via").(string)

	return res
}

func parseEgressRoute(resp *http.Response) (*EgressRoute, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	e := &EgressRoute{}
	err = json.Unmarshal(body, e)
	if err != nil {
		return nil, fmt.Errorf("could not parse egress route response: %v", err)
	}
	return e, nil
}

func CreateEgressRoute(ctx context.Context, c *Client, e *EgressRoute) (*EgressRoute, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, egressRouteEndpoint)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("could not convert egress route to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseEgressRoute(resp)
}

func GetEgressRoute(ctx context.Context, c *Client, eID string) (*EgressRoute, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, egressRouteEndpoint, eID)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseEgressRoute(resp)
}

func UpdateEgressRoute(ctx context.Context, c *Client, eID string, e *EgressRoute) (*EgressRoute, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, egressRouteEndpoint, eID)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("could not convert egress route to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseEgressRoute(resp)
}

func DeleteEgressRoute(ctx context.Context, c *Client, mID string) (*EgressRoute, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, egressRouteEndpoint, mID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseEgressRoute(resp)
}
