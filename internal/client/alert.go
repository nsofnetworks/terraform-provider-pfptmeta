package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const alertEndpoint string = "v1/alerts"

type SpikeCondition struct {
	MinHits    int    `json:"min_hits"`
	SpikeRatio int    `json:"spike_ratio"`
	SpikeType  string `json:"spike_type"`
	TimeDiff   int    `json:"time_diff"`
}

func newSpikeCondition(d *schema.ResourceData) *SpikeCondition {
	sCond := d.Get("spike_condition").([]interface{})
	if len(sCond) == 0 {
		return nil
	}
	s := sCond[0].(map[string]interface{})
	minHits := s["min_hits"].(int)
	sRatio := s["spike_ratio"].(int)
	sType := s["spike_type"].(string)
	timeDiff := s["time_diff"].(int)
	return &SpikeCondition{MinHits: minHits, SpikeRatio: sRatio, SpikeType: sType, TimeDiff: timeDiff}
}

type ThresholdCondition struct {
	Formula   string `json:"formula"`
	Op        string `json:"op"`
	Threshold int    `json:"threshold"`
}

func newThresholdCondition(d *schema.ResourceData) *ThresholdCondition {
	tCond := d.Get("threshold_condition").([]interface{})
	if len(tCond) == 0 {
		return nil
	}
	s := tCond[0].(map[string]interface{})
	f := s["formula"].(string)
	o := s["op"].(string)
	t := s["threshold"].(int)
	return &ThresholdCondition{Formula: f, Op: o, Threshold: t}
}

type Alert struct {
	ID                 string              `json:"id,omitempty"`
	Name               string              `json:"name,omitempty"`
	Description        string              `json:"description"`
	Channels           []string            `json:"channels"`
	Enabled            bool                `json:"enabled"`
	GroupBy            *string             `json:"group_by"`
	NotifyMessage      string              `json:"notify_message"`
	QueryText          string              `json:"query_text"`
	SourceType         string              `json:"source_type"`
	SpikeCondition     *SpikeCondition     `json:"spike_condition,omitempty"`
	ThresholdCondition *ThresholdCondition `json:"threshold_condition,omitempty"`
	Type               string              `json:"type,omitempty"`
	Window             int                 `json:"window"`
}

func NewAlert(d *schema.ResourceData) *Alert {
	res := &Alert{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Channels = ConfigToStringSlice("channels", d)
	res.Enabled = d.Get("enabled").(bool)
	groupBy := d.Get("group_by").(string)
	if groupBy == "" {
		res.GroupBy = nil
	} else {
		res.GroupBy = &groupBy
	}
	res.NotifyMessage = d.Get("notify_message").(string)
	res.QueryText = d.Get("query_text").(string)
	res.SourceType = d.Get("source_type").(string)
	res.SpikeCondition = newSpikeCondition(d)
	res.ThresholdCondition = newThresholdCondition(d)
	res.Window = d.Get("window").(int)
	return res
}

func parseAlert(resp []byte) (*Alert, error) {
	a := &Alert{}
	err := json.Unmarshal(resp, a)
	if err != nil {
		return nil, fmt.Errorf("could not parse alert response: %v", err)
	}
	return a, nil
}

func CreateAlert(ctx context.Context, c *Client, a *Alert) (*Alert, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, alertEndpoint)
	body, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("could not convert alert to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseAlert(resp)
}

func UpdateAlert(ctx context.Context, c *Client, aID string, a *Alert) (*Alert, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, alertEndpoint, aID)
	body, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("could not convert alert to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseAlert(resp)
}

func GetAlert(ctx context.Context, c *Client, aID string) (*Alert, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, alertEndpoint, aID)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseAlert(resp)
}

func DeleteAlert(ctx context.Context, c *Client, aID string) (*Alert, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, alertEndpoint, aID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseAlert(resp)
}
