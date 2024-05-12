package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	u "net/url"
)

const aacRuleEndpoint = "v1/aac_rules"

type AacRule struct {
	ID                   string    `json:"id,omitempty"`
	Name                 string    `json:"name,omitempty"`
	Description          *string   `json:"description"`
	Enabled              bool      `json:"enabled"`
	Priority             int       `json:"priority,omitempty"`
	Action               string    `json:"action,omitempty"`
	AppIds               []string  `json:"app_ids,omitempty"`
	ApplyAllApps         bool      `json:"apply_all_apps"`
	Sources              []string  `json:"sources,omitempty"`
	ExemptSources        []string  `json:"exempt_sources,omitempty"`
	FilterExpression     *string   `json:"filter_expression"`
	Networks             []string  `json:"networks,omitempty"`
	Locations            *[]string `json:"locations"`
	IpReputations        *[]string `json:"ip_reputations"`
	CertificateId        *string   `json:"certificate_id"`
	NotificationChannels []string  `json:"notification_channels,omitempty"`
}

func NewAacRule(d *schema.ResourceData) *AacRule {
	res := &AacRule{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Enabled = d.Get("enabled").(bool)
	res.Priority = d.Get("priority").(int)
	res.Action = d.Get("action").(string)
	res.AppIds = ResourceTypeSetToStringSlice(d.Get("app_ids").(*schema.Set))
	res.ApplyAllApps = d.Get("apply_all_apps").(bool)
	res.Sources = ResourceTypeSetToStringSlice(d.Get("sources").(*schema.Set))
	res.ExemptSources = ResourceTypeSetToStringSlice(d.Get("exempt_sources").(*schema.Set))
	res.Networks = ResourceTypeSetToStringSlice(d.Get("networks").(*schema.Set))
	res.NotificationChannels = ResourceTypeSetToStringSlice(d.Get("notification_channels").(*schema.Set))

	dsc := d.Get("description").(string)
	if dsc != "" {
		res.Description = &dsc
	} else {
		res.Description = nil
	}
	loc := ResourceTypeSetToStringSlice(d.Get("locations").(*schema.Set))
	if len(loc) > 0 {
		res.Locations = &loc
	} else {
		res.Locations = nil
	}
	ip_rep := ResourceTypeSetToStringSlice(d.Get("ip_reputations").(*schema.Set))
	if len(ip_rep) > 0 {
		res.IpReputations = &ip_rep
	} else {
		res.IpReputations = nil
	}
	cert_id := d.Get("certificate_id").(string)
	if cert_id != "" {
		res.CertificateId = &cert_id
	} else {
		res.CertificateId = nil
	}
	fe := d.Get("filter_expression").(string)
	if fe != "" {
		res.FilterExpression = &fe
	} else {
		res.FilterExpression = nil
	}
	return res
}

func parseAacRule(resp []byte) (*AacRule, error) {
	aac_rule := &AacRule{}
	err := json.Unmarshal(resp, aac_rule)
	if err != nil {
		return nil, fmt.Errorf("could not parse aac rule response: %v", err)
	}
	return aac_rule, nil
}

func CreateAacRule(ctx context.Context, c *Client, aac_rule *AacRule) (*AacRule, error) {
	arlUrl := fmt.Sprintf("%s/%s", c.BaseURL, aacRuleEndpoint)
	body, err := json.Marshal(aac_rule)
	if err != nil {
		return nil, fmt.Errorf("could not convert aac rule to json: %v", err)
	}
	resp, err := c.Post(ctx, arlUrl, body)
	if err != nil {
		return nil, err
	}
	return parseAacRule(resp)
}

func UpdateAacRule(ctx context.Context, c *Client, arlID string, aac_rule *AacRule) (*AacRule, error) {
	arlUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, aacRuleEndpoint, arlID)
	body, err := json.Marshal(aac_rule)
	if err != nil {
		return nil, fmt.Errorf("could not convert aac rule to json: %v", err)
	}
	resp, err := c.Patch(ctx, arlUrl, body)
	if err != nil {
		return nil, err
	}
	return parseAacRule(resp)
}

func GetAacRule(ctx context.Context, c *Client, arlID string) (*AacRule, error) {
	arlUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, aacRuleEndpoint, arlID)
	resp, err := c.Get(ctx, arlUrl, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseAacRule(resp)
}

func DeleteAacRule(ctx context.Context, c *Client, arlID string) (*AacRule, error) {
	arlUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, aacRuleEndpoint, arlID)
	resp, err := c.Delete(ctx, arlUrl, nil)
	if err != nil {
		return nil, err
	}
	return parseAacRule(resp)
}
