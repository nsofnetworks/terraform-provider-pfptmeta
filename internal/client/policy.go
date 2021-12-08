package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"net/http"
)

const policyEndpoint string = "v1/policies"

type Policy struct {
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description"`
	Destinations   []string `json:"destinations"`
	Enabled        *bool    `json:"enabled,omitempty"`
	ExemptSources  []string `json:"exempt_sources"`
	ProtocolGroups []string `json:"protocol_groups"`
	Sources        []string `json:"sources"`
}

func NewPolicy(d *schema.ResourceData) *Policy {
	res := &Policy{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)

	enabled := d.Get("enabled").(bool)
	res.Enabled = &enabled

	res.Destinations = ResourceTypeSetToStringSlice(d.Get("destinations").(*schema.Set))

	res.ExemptSources = ResourceTypeSetToStringSlice(d.Get("exempt_sources").(*schema.Set))

	res.Sources = ResourceTypeSetToStringSlice(d.Get("sources").(*schema.Set))

	res.ProtocolGroups = ResourceTypeSetToStringSlice(d.Get("protocol_groups").(*schema.Set))
	return res
}

func parsePolicy(resp *http.Response) (*Policy, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	pg := &Policy{}
	err = json.Unmarshal(body, pg)
	if err != nil {
		return nil, fmt.Errorf("could not parse policy response: %v", err)
	}
	return pg, nil
}

func CreatePolicy(c *Client, rg *Policy) (*Policy, error) {
	rgUrl := fmt.Sprintf("%s/%s", c.BaseURL, policyEndpoint)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert policy to json: %v", err)
	}
	resp, err := c.Post(rgUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parsePolicy(resp)
}

func UpdatePolicy(c *Client, rgID string, rg *Policy) (*Policy, error) {
	rgUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, policyEndpoint, rgID)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert policy to json: %v", err)
	}
	resp, err := c.Patch(rgUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parsePolicy(resp)
}

func GetPolicy(c *Client, rgID string) (*Policy, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, policyEndpoint, rgID)
	resp, err := c.Get(url, nil)
	if err != nil {
		return nil, err
	}
	return parsePolicy(resp)
}

func DeletePolicy(c *Client, pgID string) (*Policy, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, policyEndpoint, pgID)
	resp, err := c.Delete(url, nil)
	if err != nil {
		return nil, err
	}
	return parsePolicy(resp)
}
