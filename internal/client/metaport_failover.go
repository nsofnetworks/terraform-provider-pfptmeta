package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	u "net/url"
)

const metaportFailoverEndpoint string = "v1/metaport_failovers"

type FailBack struct {
	Trigger string `json:"trigger"`
}

type FailOver struct {
	Delay     uint8  `json:"delay"`
	Threshold uint8  `json:"threshold"`
	Trigger   string `json:"trigger"`
}

type MetaportFailover struct {
	ID                   string    `json:"id,omitempty"`
	Name                 string    `json:"name,omitempty"`
	Description          string    `json:"description,omitempty"`
	MappedElements       []string  `json:"mapped_elements"`
	Cluster1             *string   `json:"cluster_1"`
	Cluster2             *string   `json:"cluster_2"`
	ActiveCluster        *string   `json:"active_cluster,omitempty"`
	FailBack             *FailBack `json:"failback,omitempty"`
	FailOver             *FailOver `json:"failover,omitempty"`
	NotificationChannels []string  `json:"notification_channels"`
}

func (mf *MetaportFailover) ReqBody() ([]byte, error) {
	if *mf.Cluster1 == "" {
		mf.Cluster1 = nil
	}
	if *mf.Cluster2 == "" {
		mf.Cluster2 = nil
	}
	return json.Marshal(mf)
}

func NewMetaportFailover(d *schema.ResourceData) *MetaportFailover {
	res := &MetaportFailover{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)

	mes := d.Get("mapped_elements")
	res.MappedElements = ResourceTypeSetToStringSlice(mes.(*schema.Set))

	cluster1 := d.Get("cluster_1").(string)
	res.Cluster1 = &cluster1

	cluster2 := d.Get("cluster_2").(string)
	res.Cluster2 = &cluster2

	if f, exists := d.GetOk("failback"); exists {
		failback := f.([]interface{})
		if len(failback) == 1 {
			t := failback[0].(map[string]interface{})
			trigger := t["trigger"].(string)
			res.FailBack = &FailBack{Trigger: trigger}
		}
	}

	if f, exists := d.GetOk("failover"); exists {
		failover := f.([]interface{})
		if len(failover) == 1 {
			failover := failover[0].(map[string]interface{})
			delay := failover["delay"].(int)
			threshold := failover["threshold"].(int)
			trigger := failover["trigger"].(string)
			res.FailOver = &FailOver{Delay: uint8(delay), Threshold: uint8(threshold), Trigger: trigger}
		}
	}
	res.NotificationChannels = ConfigToStringSlice("notification_channels", d)

	return res
}

func parseMetaportFailover(resp []byte) (*MetaportFailover, error) {
	mf := &MetaportFailover{}
	err := json.Unmarshal(resp, mf)
	if err != nil {
		return nil, fmt.Errorf("could not parse metaport failover response: %v", err)
	}
	return mf, nil
}

func CreateMetaportFailover(ctx context.Context, c *Client, m *MetaportFailover) (*MetaportFailover, error) {
	neUrl := fmt.Sprintf("%s/%s", c.BaseURL, metaportFailoverEndpoint)
	body, err := m.ReqBody()
	if err != nil {
		return nil, fmt.Errorf("could not convert metaport failover to json: %v", err)
	}
	resp, err := c.Post(ctx, neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseMetaportFailover(resp)
}

func GetMetaportFailover(ctx context.Context, c *Client, mId string) (*MetaportFailover, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, metaportFailoverEndpoint, mId)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseMetaportFailover(resp)
}

func UpdateMetaportFailover(ctx context.Context, c *Client, mId string, m *MetaportFailover) (*MetaportFailover, error) {
	neUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, metaportFailoverEndpoint, mId)
	body, err := m.ReqBody()
	if err != nil {
		return nil, fmt.Errorf("could not convert metaport failover to json: %v", err)
	}
	resp, err := c.Patch(ctx, neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseMetaportFailover(resp)
}

func DeleteMetaportFailover(ctx context.Context, c *Client, mcID string) (*MetaportFailover, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, metaportFailoverEndpoint, mcID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseMetaportFailover(resp)
}
