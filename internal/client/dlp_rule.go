package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"net/http"
	u "net/url"
)

const dlpRulesEndpoint string = "v1/dlp_rules"

type DLPRule struct {
	ID                    string   `json:"id,omitempty"`
	Name                  string   `json:"name,omitempty"`
	Description           string   `json:"description"`
	ApplyToOrg            bool     `json:"apply_to_org"`
	Sources               []string `json:"sources,omitempty"`
	ExemptSources         []string `json:"exempt_sources,omitempty"`
	Enabled               bool     `json:"enabled"`
	Action                string   `json:"action"`
	AlertLevel            string   `json:"alert_level,omitempty"`
	AllResources          bool     `json:"all_resources"`
	AllSupportedFileTypes bool     `json:"all_supported_file_types"`
	CloudApps             []string `json:"cloud_apps"`
	ContentTypes          []string `json:"content_types"`
	Detectors             []string `json:"detectors,omitempty"`
	FileParts             []string `json:"file_parts"`
	FileTypes             []string `json:"file_types"`
	FilterExpression      *string  `json:"filter_expression"`
	Priority              int      `json:"priority"`
	ResourceCountries     []string `json:"resource_countries"`
	ThreatTypes           []string `json:"threat_types"`
	UserActions           []string `json:"user_actions"`
}

func NewDLPRule(d *schema.ResourceData) *DLPRule {
	res := &DLPRule{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Action = d.Get("action").(string)
	res.ApplyToOrg = d.Get("apply_to_org").(bool)
	res.Enabled = d.Get("enabled").(bool)
	res.Sources = ConfigToStringSlice("sources", d)
	res.ExemptSources = ConfigToStringSlice("exempt_sources", d)
	res.AlertLevel = d.Get("alert_level").(string)
	res.AllResources = d.Get("all_resources").(bool)
	res.AllSupportedFileTypes = d.Get("all_supported_file_types").(bool)
	res.CloudApps = ConfigToStringSlice("cloud_apps", d)
	res.ContentTypes = ConfigToStringSlice("content_types", d)
	res.Detectors = ConfigToStringSlice("detectors", d)
	res.FileParts = ConfigToStringSlice("file_parts", d)
	res.FileTypes = ConfigToStringSlice("file_types", d)
	fExpression := d.Get("filter_expression").(string)
	if fExpression == "" {
		res.FilterExpression = nil
	} else {
		res.FilterExpression = &fExpression
	}
	res.Priority = d.Get("priority").(int)
	res.ResourceCountries = ConfigToStringSlice("resource_countries", d)
	res.ThreatTypes = ConfigToStringSlice("threat_types", d)
	res.UserActions = ConfigToStringSlice("user_actions", d)

	return res
}

func parseDLPRule(resp *http.Response) (*DLPRule, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	pg := &DLPRule{}
	err = json.Unmarshal(body, pg)
	if err != nil {
		return nil, fmt.Errorf("could not parse dlp rule response: %v", err)
	}
	return pg, nil
}

func CreateDLPRule(ctx context.Context, c *Client, rg *DLPRule) (*DLPRule, error) {
	rgUrl := fmt.Sprintf("%s/%s", c.BaseURL, dlpRulesEndpoint)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert dlp rule to json: %v", err)
	}
	resp, err := c.Post(ctx, rgUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseDLPRule(resp)
}

func UpdateDLPRule(ctx context.Context, c *Client, rgID string, rg *DLPRule) (*DLPRule, error) {
	rgUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, dlpRulesEndpoint, rgID)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert dlp rule to json: %v", err)
	}
	resp, err := c.Patch(ctx, rgUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseDLPRule(resp)
}

func GetDLPRule(ctx context.Context, c *Client, rgID string) (*DLPRule, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, dlpRulesEndpoint, rgID)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseDLPRule(resp)
}

func DeleteDLPRule(ctx context.Context, c *Client, pgID string) (*DLPRule, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, dlpRulesEndpoint, pgID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseDLPRule(resp)
}
