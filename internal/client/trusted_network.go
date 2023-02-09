package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const trustedNetworkEndpoint = "v1/trusted_networks"

type ExternalIpConfig struct {
	AddressesRanges []string `json:"addresses_ranges"`
}

func newExternalIPConfig(input interface{}) *ExternalIpConfig {
	newExternalIpConfig := &ExternalIpConfig{}
	inputList := input.([]interface{})
	if len(inputList) == 0 {
		return nil
	}
	externalIpConfig := inputList[0].(map[string]interface{})
	addressesRanges := externalIpConfig["addresses_ranges"].([]interface{})
	newExternalIpConfig.AddressesRanges = make([]string, len(addressesRanges))
	for j, address := range addressesRanges {
		newExternalIpConfig.AddressesRanges[j] = address.(string)
	}
	return newExternalIpConfig
}

type ResolvedAddressConfig struct {
	AddressesRanges []string `json:"addresses_ranges"`
	Hostname        string   `json:"hostname"`
}

func newResolvedAddressConfig(input interface{}) *ResolvedAddressConfig {
	newResolvedAddressConfig := &ResolvedAddressConfig{}
	inputList := input.([]interface{})
	if len(inputList) == 0 {
		return nil
	}
	resolvedAddressConfig := inputList[0].(map[string]interface{})
	addressesRanges := resolvedAddressConfig["addresses_ranges"].([]interface{})
	newResolvedAddressConfig.AddressesRanges = make([]string, len(addressesRanges))
	for j, address := range addressesRanges {
		newResolvedAddressConfig.AddressesRanges[j] = address.(string)
	}
	newResolvedAddressConfig.Hostname = resolvedAddressConfig["hostname"].(string)
	return newResolvedAddressConfig
}

type Criteria struct {
	ExternalIpConfig      *ExternalIpConfig      `json:"external_ip_config,omitempty"`
	ResolvedAddressConfig *ResolvedAddressConfig `json:"resolved_address_config,omitempty"`
	Type                  string                 `json:"type,omitempty"`
}

func newCriteria(d *schema.ResourceData) []Criteria {
	c := d.Get("criteria").([]interface{})
	res := make([]Criteria, len(c))
	for i, criteria := range c {
		newCriteria := Criteria{}
		criteria := criteria.(map[string]interface{})
		if externalIpConfig, ok := criteria["external_ip_config"]; ok {
			newCriteria.ExternalIpConfig = newExternalIPConfig(externalIpConfig)
		}
		if resolvedAddressConfig, ok := criteria["resolved_address_config"]; ok {
			newCriteria.ResolvedAddressConfig = newResolvedAddressConfig(resolvedAddressConfig)
		}
		res[i] = newCriteria
	}
	return res
}

type TrustedNetwork struct {
	ID              string     `json:"id,omitempty"`
	Name            string     `json:"name,omitempty"`
	Description     string     `json:"description"`
	Enabled         bool       `json:"enabled"`
	ApplyToOrg      bool       `json:"apply_to_org"`
	ApplyToEntities []string   `json:"apply_to_entities"`
	ExemptEntities  []string   `json:"exempt_entities"`
	Criteria        []Criteria `json:"criteria"`
}

func NewTrustedNetwork(d *schema.ResourceData) *TrustedNetwork {
	res := &TrustedNetwork{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Enabled = d.Get("enabled").(bool)
	res.ApplyToOrg = d.Get("apply_to_org").(bool)
	res.ApplyToEntities = ConfigToStringSlice("apply_to_entities", d)
	res.ExemptEntities = ConfigToStringSlice("exempt_entities", d)
	res.Criteria = newCriteria(d)
	return res
}

func parseTrustedNetwork(resp []byte) (*TrustedNetwork, error) {
	e := &TrustedNetwork{}
	err := json.Unmarshal(resp, e)
	if err != nil {
		return nil, fmt.Errorf("could not parse trusted network response: %v", err)
	}
	return e, nil
}

func CreateTrustedNetwork(ctx context.Context, c *Client, e *TrustedNetwork) (*TrustedNetwork, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, trustedNetworkEndpoint)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("could not convert trusted network to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseTrustedNetwork(resp)
}

func GetTrustedNetwork(ctx context.Context, c *Client, eID string) (*TrustedNetwork, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, trustedNetworkEndpoint, eID)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseTrustedNetwork(resp)
}

func UpdateTrustedNetwork(ctx context.Context, c *Client, eID string, e *TrustedNetwork) (*TrustedNetwork, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, trustedNetworkEndpoint, eID)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("could not convert trusted network to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseTrustedNetwork(resp)
}

func DeleteTrustedNetwork(ctx context.Context, c *Client, mID string) (*TrustedNetwork, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, trustedNetworkEndpoint, mID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseTrustedNetwork(resp)
}
