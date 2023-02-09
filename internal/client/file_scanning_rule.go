package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	u "net/url"
)

const fileScanningRulesEndpoint string = "v1/file_scanning_rules"

type FileScanningRule struct {
	ID                    string   `json:"id,omitempty"`
	Name                  string   `json:"name,omitempty"`
	Description           string   `json:"description"`
	ApplyToOrg            bool     `json:"apply_to_org"`
	Enabled               bool     `json:"enabled"`
	Sources               []string `json:"sources,omitempty"`
	ExemptSources         []string `json:"exempt_sources,omitempty"`
	CloudApps             []string `json:"cloud_apps,omitempty"`
	BlockAllFileTypes     bool     `json:"block_all_file_types"`
	BlockContentTypes     []string `json:"block_content_types,omitempty"`
	BlockCountries        []string `json:"block_countries,omitempty"`
	BlockFileTypes        []string `json:"block_file_types,omitempty"`
	BlockThreatTypes      []string `json:"block_threat_types,omitempty"`
	BlockUnsupportedFiles bool     `json:"block_unsupported_files"`
	FilterExpression      *string  `json:"filter_expression"`
	Malware               string   `json:"malware"`
	MaxFileSizeMb         int      `json:"max_file_size_mb,omitempty"`
	Priority              int      `json:"priority"`
	SandboxFileTypes      []string `json:"sandbox_file_types,omitempty"`
	TimeoutPolicy         string   `json:"timeout_policy,omitempty"`
}

func NewFileScanningRule(d *schema.ResourceData) *FileScanningRule {
	res := &FileScanningRule{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.ApplyToOrg = d.Get("apply_to_org").(bool)
	res.Enabled = d.Get("enabled").(bool)
	res.Sources = ConfigToStringSlice("sources", d)
	res.ExemptSources = ConfigToStringSlice("exempt_sources", d)
	res.BlockAllFileTypes = d.Get("block_all_file_types").(bool)
	res.BlockContentTypes = ConfigToStringSlice("block_content_types", d)
	res.CloudApps = ConfigToStringSlice("cloud_apps", d)
	res.BlockCountries = ConfigToStringSlice("block_countries", d)
	res.BlockFileTypes = ConfigToStringSlice("block_file_types", d)
	res.BlockThreatTypes = ConfigToStringSlice("block_threat_types", d)
	res.BlockUnsupportedFiles = d.Get("block_unsupported_files").(bool)
	fExpression := d.Get("filter_expression").(string)
	if fExpression == "" {
		res.FilterExpression = nil
	} else {
		res.FilterExpression = &fExpression
	}
	res.Malware = d.Get("malware").(string)
	res.MaxFileSizeMb = d.Get("max_file_size_mb").(int)
	res.Priority = d.Get("priority").(int)
	res.SandboxFileTypes = ConfigToStringSlice("sandbox_file_types", d)
	res.TimeoutPolicy = d.Get("timeout_policy").(string)

	return res
}

func parseFileScanningRule(resp []byte) (*FileScanningRule, error) {
	pg := &FileScanningRule{}
	err := json.Unmarshal(resp, pg)
	if err != nil {
		return nil, fmt.Errorf("could not parse file scanning rule response: %v", err)
	}
	return pg, nil
}

func CreateFileScanningRule(ctx context.Context, c *Client, rg *FileScanningRule) (*FileScanningRule, error) {
	rgUrl := fmt.Sprintf("%s/%s", c.BaseURL, fileScanningRulesEndpoint)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert file scanning rule to json: %v", err)
	}
	resp, err := c.Post(ctx, rgUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseFileScanningRule(resp)
}

func UpdateFileScanningRule(ctx context.Context, c *Client, rgID string, rg *FileScanningRule) (*FileScanningRule, error) {
	rgUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, fileScanningRulesEndpoint, rgID)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert file scanning rule to json: %v", err)
	}
	resp, err := c.Patch(ctx, rgUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseFileScanningRule(resp)
}

func GetFileScanningRule(ctx context.Context, c *Client, rgID string) (*FileScanningRule, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, fileScanningRulesEndpoint, rgID)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseFileScanningRule(resp)
}

func DeleteFileScanningRule(ctx context.Context, c *Client, pgID string) (*FileScanningRule, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, fileScanningRulesEndpoint, pgID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseFileScanningRule(resp)
}
