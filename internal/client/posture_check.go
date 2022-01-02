package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"net/http"
)

type Check struct {
	MinVersion string `json:"min_version,omitempty"`
	Type       string `json:"type"`
}

func newCheck(d *schema.ResourceData) *Check {
	check := d.Get("check").([]interface{})
	if len(check) == 1 {
		check := check[0].(map[string]interface{})
		return &Check{Type: check["type"].(string), MinVersion: check["min_version"].(string)}
	}
	return nil
}

type PostureCheck struct {
	ID                string   `json:"id,omitempty"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Enabled           bool     `json:"enabled"`
	Action            string   `json:"action"`
	Check             *Check   `json:"check,omitempty"`
	When              []string `json:"when"`
	ApplyToOrg        bool     `json:"apply_to_org"`
	ApplyToEntities   []string `json:"apply_to_entities"`
	ExemptEntities    []string `json:"exempt_entities"`
	Interval          *int     `json:"interval"`
	Osquery           string   `json:"osquery"`
	Platform          string   `json:"platform"`
	UserMessageOnFail string   `json:"user_message_on_fail"`
}

const postureCheckEndpoint = "v1/posture_checks"

func NewPostureCheck(d *schema.ResourceData) *PostureCheck {
	res := &PostureCheck{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Enabled = d.Get("enabled").(bool)
	res.Action = d.Get("action").(string)
	res.Check = newCheck(d)
	when := d.Get("when").([]interface{})
	res.When = make([]string, len(when))
	for i, val := range when {
		res.When[i] = val.(string)
	}
	res.ApplyToOrg = d.Get("apply_to_org").(bool)
	applyToEntities := d.Get("apply_to_entities").([]interface{})
	res.ApplyToEntities = make([]string, len(applyToEntities))
	for i, val := range applyToEntities {
		res.ApplyToEntities[i] = val.(string)
	}
	exemptEntities := d.Get("exempt_entities").([]interface{})
	res.ExemptEntities = make([]string, len(exemptEntities))
	for i, val := range exemptEntities {
		res.ExemptEntities[i] = val.(string)
	}
	interval := d.Get("interval").(int)
	if interval == 0 {
		res.Interval = nil
	} else {
		res.Interval = &interval
	}
	res.Osquery = d.Get("osquery").(string)
	res.Platform = d.Get("platform").(string)
	res.UserMessageOnFail = d.Get("user_message_on_fail").(string)
	return res
}

func parsePostureCheck(resp *http.Response) (*PostureCheck, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	e := &PostureCheck{}
	err = json.Unmarshal(body, e)
	if err != nil {
		return nil, fmt.Errorf("could not parse posture check response: %v", err)
	}
	return e, nil
}

func CreatePostureCheck(ctx context.Context, c *Client, e *PostureCheck) (*PostureCheck, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, postureCheckEndpoint)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("could not convert posture check to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parsePostureCheck(resp)
}

func GetPostureCheck(ctx context.Context, c *Client, eID string) (*PostureCheck, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, postureCheckEndpoint, eID)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parsePostureCheck(resp)
}

func UpdatePostureCheck(ctx context.Context, c *Client, eID string, e *PostureCheck) (*PostureCheck, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, postureCheckEndpoint, eID)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("could not convert posture check to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parsePostureCheck(resp)
}

func DeletePostureCheck(ctx context.Context, c *Client, mID string) (*PostureCheck, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, postureCheckEndpoint, mID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parsePostureCheck(resp)
}
