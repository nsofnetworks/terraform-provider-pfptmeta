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

const sslBypassRulesEndpoint = "v1/ssl_bypass_rules"

type SSLBypassRule struct {
	ID                      string   `json:"id,omitempty"`
	Name                    string   `json:"name"`
	Description             string   `json:"description,omitempty"`
	Enabled                 bool     `json:"enabled"`
	ApplyToOrg              bool     `json:"apply_to_org"`
	Sources                 []string `json:"sources"`
	ExemptSources           []string `json:"exempt_sources"`
	BypassUncategorizedUrls bool     `json:"bypass_uncategorized_urls"`
	ContentTypes            []string `json:"content_types"`
	Domains                 []string `json:"domains"`
	Priority                int      `json:"priority"`
}

func NewSSLBypassRule(d *schema.ResourceData) *SSLBypassRule {
	res := &SSLBypassRule{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.ApplyToOrg = d.Get("apply_to_org").(bool)
	res.Enabled = d.Get("enabled").(bool)
	res.Sources = ConfigToStringSlice("sources", d)
	res.ExemptSources = ConfigToStringSlice("exempt_sources", d)
	res.Priority = d.Get("priority").(int)
	res.Domains = ConfigToStringSlice("domains", d)
	res.ContentTypes = ConfigToStringSlice("content_types", d)
	res.BypassUncategorizedUrls = d.Get("bypass_uncategorized_urls").(bool)

	return res
}

func parseSSLBypassRule(resp *http.Response) (*SSLBypassRule, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	pg := &SSLBypassRule{}
	err = json.Unmarshal(body, pg)
	if err != nil {
		return nil, fmt.Errorf("could not parse ssl bypass rule response: %v", err)
	}
	return pg, nil
}

func CreateSSLBypassRule(ctx context.Context, c *Client, rg *SSLBypassRule) (*SSLBypassRule, error) {
	rgUrl := fmt.Sprintf("%s/%s", c.BaseURL, sslBypassRulesEndpoint)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert ssl bypass rule to json: %v", err)
	}
	resp, err := c.Post(ctx, rgUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseSSLBypassRule(resp)
}

func UpdateSSLBypassRule(ctx context.Context, c *Client, rgID string, rg *SSLBypassRule) (*SSLBypassRule, error) {
	rgUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, sslBypassRulesEndpoint, rgID)
	body, err := json.Marshal(rg)
	if err != nil {
		return nil, fmt.Errorf("could not convert ssl bypass rule to json: %v", err)
	}
	resp, err := c.Patch(ctx, rgUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseSSLBypassRule(resp)
}

func GetSSLBypassRule(ctx context.Context, c *Client, rgID string) (*SSLBypassRule, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, sslBypassRulesEndpoint, rgID)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseSSLBypassRule(resp)
}

func DeleteSSLBypassRule(ctx context.Context, c *Client, pgID string) (*SSLBypassRule, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, sslBypassRulesEndpoint, pgID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseSSLBypassRule(resp)
}
