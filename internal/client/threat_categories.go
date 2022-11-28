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

const threatCategoryEndpoint = "v1/threat_categories"

type ThreatCategory struct {
	ID              string   `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description"`
	ConfidenceLevel string   `json:"confidence_level"`
	RiskLevel       string   `json:"risk_level"`
	Countries       []string `json:"countries"`
	Types           []string `json:"types"`
}

func NewThreatCategory(d *schema.ResourceData) *ThreatCategory {
	res := &ThreatCategory{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.ConfidenceLevel = d.Get("confidence_level").(string)
	res.RiskLevel = d.Get("risk_level").(string)
	res.Countries = ConfigToStringSlice("countries", d)
	res.Types = ConfigToStringSlice("types", d)
	return res
}

func parseThreatCategory(resp *http.Response) (*ThreatCategory, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read threat category response: %v", err)
	}
	tc := &ThreatCategory{}
	err = json.Unmarshal(body, tc)
	if err != nil {
		return nil, fmt.Errorf("could not parse threat category response: %v", err)
	}
	return tc, nil
}

func CreateThreatCategory(ctx context.Context, c *Client, tc *ThreatCategory) (*ThreatCategory, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, threatCategoryEndpoint)
	body, err := json.Marshal(tc)
	if err != nil {
		return nil, fmt.Errorf("could not convert threat category to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseThreatCategory(resp)
}

func UpdateThreatCategory(ctx context.Context, c *Client, tcId string, cc *ThreatCategory) (*ThreatCategory, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, threatCategoryEndpoint, tcId)
	body, err := json.Marshal(cc)
	if err != nil {
		return nil, fmt.Errorf("could not convert threat category to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseThreatCategory(resp)
}

func GetThreatCategory(ctx context.Context, c *Client, tcId string) (*ThreatCategory, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, threatCategoryEndpoint, tcId)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseThreatCategory(resp)
}

func DeleteThreatCategory(ctx context.Context, c *Client, tcId string) (*ThreatCategory, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, threatCategoryEndpoint, tcId)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseThreatCategory(resp)
}
