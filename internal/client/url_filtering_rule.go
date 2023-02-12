package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	u "net/url"
)

const urlFilteringRulesEndpoint string = "v1/url_filtering_rules"

type UrlFilteringRule struct {
	ID                         string   `json:"id,omitempty"`
	Name                       string   `json:"name,omitempty"`
	Description                string   `json:"description"`
	Action                     string   `json:"action"`
	ApplyToOrg                 bool     `json:"apply_to_org"`
	Enabled                    bool     `json:"enabled"`
	Sources                    []string `json:"sources,omitempty"`
	ExemptSources              []string `json:"exempt_sources,omitempty"`
	AdvancedThreatProtection   bool     `json:"advanced_threat_protection"`
	CatalogAppCategories       []string `json:"catalog_app_categories"`
	CatalogAppRisk             int      `json:"catalog_app_risk,omitempty"`
	CloudApps                  []string `json:"cloud_apps"`
	Countries                  []string `json:"countries,omitempty"`
	ExpiresAt                  string   `json:"expires_at,omitempty"`
	FilterExpression           string   `json:"filter_expression,omitempty"`
	ForbiddenContentCategories []string `json:"forbidden_content_categories"`
	Networks                   []string `json:"networks,omitempty"`
	Priority                   int      `json:"priority"`
	Schedule                   []string `json:"schedule"`
	TenantRestriction          string   `json:"tenant_restriction,omitempty"`
	ThreatCategories           []string `json:"threat_categories"`
	WarnTtl                    int      `json:"warn_ttl"`
}

func NewUrlFilteringRule(d *schema.ResourceData) *UrlFilteringRule {
	res := &UrlFilteringRule{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Action = d.Get("action").(string)
	res.ApplyToOrg = d.Get("apply_to_org").(bool)
	res.Enabled = d.Get("enabled").(bool)
	res.Sources = ConfigToStringSlice("sources", d)
	res.ExemptSources = ConfigToStringSlice("exempt_sources", d)
	res.AdvancedThreatProtection = d.Get("advanced_threat_protection").(bool)
	res.CatalogAppCategories = ConfigToStringSlice("catalog_app_categories", d)
	res.CatalogAppRisk = d.Get("catalog_app_risk").(int)
	res.CloudApps = ConfigToStringSlice("cloud_apps", d)
	res.Countries = ConfigToStringSlice("countries", d)
	res.ExpiresAt = d.Get("expires_at").(string)
	res.FilterExpression = d.Get("filter_expression").(string)
	res.ForbiddenContentCategories = ConfigToStringSlice("forbidden_content_categories", d)
	res.Networks = ConfigToStringSlice("networks", d)
	res.Priority = d.Get("priority").(int)
	res.Schedule = ConfigToStringSlice("schedule", d)
	res.TenantRestriction = d.Get("tenant_restriction").(string)
	res.ThreatCategories = ConfigToStringSlice("threat_categories", d)
	res.WarnTtl = d.Get("warn_ttl").(int)

	return res
}

func parseUrlFilteringRule(resp []byte) (*UrlFilteringRule, error) {
	pg := &UrlFilteringRule{}
	err := json.Unmarshal(resp, pg)
	if err != nil {
		return nil, fmt.Errorf("could not parse url filtering rule response: %v", err)
	}
	return pg, nil
}

func CreateUrlFilteringRule(ctx context.Context, c *Client, rg *UrlFilteringRule) (*UrlFilteringRule, error) {
	rgUrl := fmt.Sprintf("%s/%s", c.BaseURL, urlFilteringRulesEndpoint)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert url filtering rule to json: %v", err)
	}
	resp, err := c.Post(ctx, rgUrl, body)
	if err != nil {
		return nil, err
	}
	return parseUrlFilteringRule(resp)
}

func UpdateUrlFilteringRule(ctx context.Context, c *Client, rgID string, rg *UrlFilteringRule) (*UrlFilteringRule, error) {
	rgUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, urlFilteringRulesEndpoint, rgID)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert url filtering rule to json: %v", err)
	}
	resp, err := c.Patch(ctx, rgUrl, body)
	if err != nil {
		return nil, err
	}
	return parseUrlFilteringRule(resp)
}

func GetUrlFilteringRule(ctx context.Context, c *Client, rgID string) (*UrlFilteringRule, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, urlFilteringRulesEndpoint, rgID)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseUrlFilteringRule(resp)
}

func DeleteUrlFilteringRule(ctx context.Context, c *Client, pgID string) (*UrlFilteringRule, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, urlFilteringRulesEndpoint, pgID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseUrlFilteringRule(resp)
}
