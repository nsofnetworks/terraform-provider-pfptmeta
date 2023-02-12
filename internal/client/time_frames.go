package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const TimeFramesEndpoint = "v1/time_frames"

type Time struct {
	Hour   int `json:"hour"`
	Minute int `json:"minute"`
}

func newTime(timeKey string, d *schema.ResourceData) *Time {
	res := &Time{}
	rawTime := d.Get(timeKey).([]interface{})
	if len(rawTime) == 0 {
		return nil
	}
	conf := rawTime[0].(map[string]interface{})
	res.Hour = conf["hour"].(int)
	res.Minute = conf["minute"].(int)
	return res
}

type TimeFrame struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description"`
	Days        []string `json:"days"`
	StartTime   *Time    `json:"start_time"`
	EndTime     *Time    `json:"end_time"`
}

func NewTimeFrame(d *schema.ResourceData) *TimeFrame {
	res := &TimeFrame{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Days = ConfigToStringSlice("days", d)
	res.StartTime = newTime("start_time", d)
	res.EndTime = newTime("end_time", d)
	return res
}

func parseTimeFrame(resp []byte) (*TimeFrame, error) {
	tf := &TimeFrame{}
	err := json.Unmarshal(resp, tf)
	if err != nil {
		return nil, fmt.Errorf("could not parse time frame response: %v", err)
	}
	return tf, nil
}

func CreateTimeFrame(ctx context.Context, c *Client, tf *TimeFrame) (*TimeFrame, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, TimeFramesEndpoint)
	body, err := json.Marshal(tf)
	if err != nil {
		return nil, fmt.Errorf("could not convert time frame to json: %v", err)
	}
	resp, err := c.Post(ctx, url, body)
	if err != nil {
		return nil, err
	}
	return parseTimeFrame(resp)
}

func UpdateTimeFrame(ctx context.Context, c *Client, tfId string, in *TimeFrame) (*TimeFrame, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, TimeFramesEndpoint, tfId)
	body, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("could not convert time frame to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, body)
	if err != nil {
		return nil, err
	}
	return parseTimeFrame(resp)
}

func GetTimeFrame(ctx context.Context, c *Client, tfId string) (*TimeFrame, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, TimeFramesEndpoint, tfId)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseTimeFrame(resp)
}

func DeleteTimeFrame(ctx context.Context, c *Client, inId string) (*TimeFrame, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, TimeFramesEndpoint, inId)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseTimeFrame(resp)
}
