package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	u "net/url"
)

const scanRulesEndpoint string = "v1/scan_rules"

type ScanRule struct {
	ID                     string   `json:"id,omitempty"`
	Priority               int      `json:"priority"`
	Name                   string   `json:"name,omitempty"`
	Description            string   `json:"description"`
	Enabled                bool     `json:"enabled"`
	Sources                []string `json:"sources,omitempty"`
	ExemptSources          []string `json:"exempt_sources,omitempty"`
	FilterExpression       *string  `json:"filter_expression,omitempty"`
	ApplyToOrg             bool     `json:"apply_to_org"`
	ContentCategories      []string `json:"content_categories,omitempty"`
	ThreatCategories       []string `json:"threat_categories,omitempty"`
	CloudApps              []string `json:"cloud_apps,omitempty"`
	CloudAppRiskGroups     []string `json:"cloud_app_risk_groups,omitempty"`
	CatalogAppCategories   []string `json:"catalog_app_categories,omitempty"`
	Networks               []string `json:"networks,omitempty"`
	Countries              []string `json:"countries,omitempty"`
	UserAgents             []string `json:"user_agents,omitempty"`
	UserActions            []string `json:"user_actions"`
	AllSupportedFileTypes  bool     `json:"all_supported_file_types"`
	FileTypes              []string `json:"file_types,omitempty"`
	MaxFileSizeMb          int      `json:"max_file_size_mb,omitempty"`
	PasswordProtectedFiles bool     `json:"password_protected_files,omitempty"`
	Dlp                    bool     `json:"dlp,omitempty"`
	Detectors              []string `json:"detectors,omitempty"`
	Malware                bool     `json:"malware,omitempty"`
	Sandbox                bool     `json:"sandbox,omitempty"`
	Antivirus              bool     `json:"antivirus,omitempty"`
	Action                 string   `json:"action"`
}

func NewScanRule(d *schema.ResourceData) *ScanRule {
	res := &ScanRule{}
	res.Priority = d.Get("priority").(int)
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Enabled = d.Get("enabled").(bool)
	res.Sources = ConfigToStringSlice("sources", d)
	res.ExemptSources = ConfigToStringSlice("exempt_sources", d)
	fExpression := d.Get("filter_expression").(string)
	if fExpression == "" {
		res.FilterExpression = nil
	} else {
		res.FilterExpression = &fExpression
	}
	res.ApplyToOrg = d.Get("apply_to_org").(bool)
	res.ContentCategories = ConfigToStringSlice("content_categories", d)
	res.ThreatCategories = ConfigToStringSlice("threat_categories", d)
	res.CloudApps = ConfigToStringSlice("cloud_apps", d)
	res.CloudAppRiskGroups = ConfigToStringSlice("cloud_app_risk_groups", d)
	res.CatalogAppCategories = ConfigToStringSlice("catalog_app_categories", d)
	res.Networks = ConfigToStringSlice("networks", d)
	res.Countries = ConfigToStringSlice("countries", d)
	res.UserAgents = ConfigToStringSlice("user_agents", d)
	res.UserActions = ConfigToStringSlice("user_actions", d)
	res.AllSupportedFileTypes = d.Get("all_supported_file_types").(bool)
	res.FileTypes = ConfigToStringSlice("file_types", d)
	res.MaxFileSizeMb = d.Get("max_file_size_mb").(int)
	res.PasswordProtectedFiles = d.Get("password_protected_files").(bool)
	res.Dlp = d.Get("dlp").(bool)
	res.Detectors = ConfigToStringSlice("detectors", d)
	res.Malware = d.Get("malware").(bool)
	res.Sandbox = d.Get("sandbox").(bool)
	res.Antivirus = d.Get("antivirus").(bool)
	res.Action = d.Get("action").(string)

	return res
}

func parseScanRule(resp []byte) (*ScanRule, error) {
	pg := &ScanRule{}
	err := json.Unmarshal(resp, pg)
	if err != nil {
		return nil, fmt.Errorf("could not parse scan rule response: %v", err)
	}
	return pg, nil
}

func CreateScanRule(ctx context.Context, c *Client, rg *ScanRule) (*ScanRule, error) {
	rgUrl := fmt.Sprintf("%s/%s", c.BaseURL, scanRulesEndpoint)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert scan rule to json: %v", err)
	}
	resp, err := c.Post(ctx, rgUrl, body)
	if err != nil {
		return nil, err
	}
	return parseScanRule(resp)
}

func UpdateScanRule(ctx context.Context, c *Client, rgID string, rg *ScanRule) (*ScanRule, error) {
	rgUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, scanRulesEndpoint, rgID)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert scan rule to json: %v", err)
	}
	resp, err := c.Patch(ctx, rgUrl, body)
	if err != nil {
		return nil, err
	}
	return parseScanRule(resp)
}

func GetScanRule(ctx context.Context, c *Client, rgID string) (*ScanRule, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, scanRulesEndpoint, rgID)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseScanRule(resp)
}

func DeleteScanRule(ctx context.Context, c *Client, pgID string) (*ScanRule, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, scanRulesEndpoint, pgID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseScanRule(resp)
}
