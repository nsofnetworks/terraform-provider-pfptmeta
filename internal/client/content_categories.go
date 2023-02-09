package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	u "net/url"
)

const contentCategoryEndpoint = "v1/content_categories"

type ContentCategory struct {
	ID                      string   `json:"id,omitempty"`
	Name                    string   `json:"name,omitempty"`
	Description             string   `json:"description"`
	ConfidenceLevel         string   `json:"confidence_level"`
	ForbidUncategorizedUrls bool     `json:"forbid_uncategorized_urls"`
	Types                   []string `json:"types"`
	Urls                    []string `json:"urls"`
}

func NewContentCategory(d *schema.ResourceData) *ContentCategory {
	res := &ContentCategory{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.ConfidenceLevel = d.Get("confidence_level").(string)
	res.ForbidUncategorizedUrls = d.Get("forbid_uncategorized_urls").(bool)
	res.Types = ConfigToStringSlice("types", d)
	res.Urls = ConfigToStringSlice("urls", d)

	return res
}

func parseContentCategory(resp []byte) (*ContentCategory, error) {
	ds := &ContentCategory{}
	err := json.Unmarshal(resp, ds)
	if err != nil {
		return nil, fmt.Errorf("could not parse content category response: %v", err)
	}
	return ds, nil
}

func CreateContentCategory(ctx context.Context, c *Client, cc *ContentCategory) (*ContentCategory, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, contentCategoryEndpoint)
	body, err := json.Marshal(cc)
	if err != nil {
		return nil, fmt.Errorf("could not convert content category to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseContentCategory(resp)
}

func UpdateContentCategory(ctx context.Context, c *Client, ccId string, cc *ContentCategory) (*ContentCategory, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, contentCategoryEndpoint, ccId)
	body, err := json.Marshal(cc)
	if err != nil {
		return nil, fmt.Errorf("could not convert content category to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseContentCategory(resp)
}

func GetContentCategory(ctx context.Context, c *Client, ccId string) (*ContentCategory, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, contentCategoryEndpoint, ccId)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseContentCategory(resp)
}

func DeleteContentCategory(ctx context.Context, c *Client, ccId string) (*ContentCategory, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, contentCategoryEndpoint, ccId)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseContentCategory(resp)
}
