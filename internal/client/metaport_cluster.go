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

const (
	metaportClusterEndpoint string = "v1/metaport_clusters"
)

type MetaportCluster struct {
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	MappedElements []string `json:"mapped_elements"`
	Metaports      []string `json:"metaports"`
}

func NewMetaportCluster(d *schema.ResourceData) *MetaportCluster {
	res := &MetaportCluster{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)

	mes := d.Get("mapped_elements")
	res.MappedElements = ResourceTypeSetToStringSlice(mes.(*schema.Set))

	mps := d.Get("metaports")
	res.Metaports = ResourceTypeSetToStringSlice(mps.(*schema.Set))

	return res
}

func parseMetaportCluster(resp *http.Response) (*MetaportCluster, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	mc := &MetaportCluster{}
	err = json.Unmarshal(body, mc)
	if err != nil {
		return nil, fmt.Errorf("could not parse metaport cluster response: %v", err)
	}
	return mc, nil
}

func CreateMetaportCluster(ctx context.Context, c *Client, m *MetaportCluster) (*MetaportCluster, error) {
	neUrl := fmt.Sprintf("%s/%s", c.BaseURL, metaportClusterEndpoint)
	body, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("could not convert metaport cluster to json: %v", err)
	}
	resp, err := c.Post(ctx, neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseMetaportCluster(resp)
}

func GetMetaportCluster(ctx context.Context, c *Client, mId string) (*MetaportCluster, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, metaportClusterEndpoint, mId)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseMetaportCluster(resp)
}

func GetMetaportClustertByName(ctx context.Context, c *Client, name string) (*MetaportCluster, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, metaportClusterEndpoint)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read metaport cluster response")
	}
	var respBody []MetaportCluster
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return nil, fmt.Errorf("could not parse metaport cluster response: %v", err)
	}
	var nameMatch []MetaportCluster
	for _, m := range respBody {
		if m.Name == name {
			nameMatch = append(nameMatch, m)
		}
	}
	switch len(nameMatch) {
	case 0:
		return nil, fmt.Errorf("could not find metaport cluster with name \"%s\"", name)
	case 1:
		return &nameMatch[0], nil
	default:
		return nil, fmt.Errorf("found more then one metaport cluster with name \"%s\"", name)
	}
}

func UpdateMetaportCluster(ctx context.Context, c *Client, mId string, m *MetaportCluster) (*MetaportCluster, error) {
	neUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, metaportClusterEndpoint, mId)
	body, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("could not convert metaport cluster to json: %v", err)
	}
	resp, err := c.Patch(ctx, neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseMetaportCluster(resp)
}

func DeleteMetaportCluster(ctx context.Context, c *Client, mcID string) (*MetaportCluster, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, metaportClusterEndpoint, mcID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseMetaportCluster(resp)
}

func AddMappedElementsToMetaportCluster(ctx context.Context, c *Client, mID string, meIDs []string) (*MetaportCluster, error) {
	url := fmt.Sprintf("%s/%s/%s/add_mapped_elements", c.BaseURL, metaportClusterEndpoint, mID)
	body := make(map[string][]string)
	body["mapped_elements"] = meIDs
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("could not convert mapped elements to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	return parseMetaportCluster(resp)
}

func RemoveMappedElementsFromMetaportCluster(ctx context.Context,
	c *Client, mID string, meIDs []string) (*MetaportCluster, error) {
	url := fmt.Sprintf("%s/%s/%s/remove_mapped_elements", c.BaseURL, metaportClusterEndpoint, mID)
	body := make(map[string][]string)
	body["mapped_elements"] = meIDs
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("could not convert mapped elements to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}
	return parseMetaportCluster(resp)
}
