package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"net/http"
)

const (
	metaportEndpoint string = "v1/metaports"
)

type Metaport struct {
	ID             string   `json:"id,omitempty"`
	Name           string   `json:"name,omitempty"`
	Description    string   `json:"description,omitempty"`
	Enabled        *bool    `json:"enabled,omitempty"`
	AllowSupport   *bool    `json:"allow_support,omitempty"`
	MappedElements []string `json:"mapped_elements"`
	Health         string   `json:"health,omitempty"`
}

func NewMetaport(d *schema.ResourceData) *Metaport {
	res := &Metaport{}
	if d.HasChange("name") {
		name := d.Get("name")
		res.Name = name.(string)
	}

	res.Description = d.Get("description").(string)

	enabled := d.Get("enabled").(bool)
	res.Enabled = &enabled

	allowSupport := d.Get("allow_support").(bool)
	res.AllowSupport = &allowSupport

	mes := d.Get("mapped_elements")
	res.MappedElements = ResourceTypeSetToStringSlice(mes.(*schema.Set))

	return res
}

func parseMetaport(resp *http.Response) (*Metaport, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	m := &Metaport{}
	err = json.Unmarshal(body, m)
	if err != nil {
		return nil, fmt.Errorf("could not parse metaport response: %v", err)
	}
	return m, nil
}

func CreateMetaport(c *Client, m *Metaport) (*Metaport, error) {
	neUrl := fmt.Sprintf("%s/%s", c.BaseURL, metaportEndpoint)
	body, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("could not convert metaport to json: %v", err)
	}
	resp, err := c.Post(neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseMetaport(resp)
}

func GetMetaport(c *Client, mId string) (*Metaport, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, metaportEndpoint, mId)
	resp, err := c.Get(url, nil)
	if err != nil {
		return nil, err
	}
	return parseMetaport(resp)
}

func UpdateMetaport(c *Client, mId string, m *Metaport) (*Metaport, error) {
	neUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, metaportEndpoint, mId)
	body, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("could not convert metaport to json: %v", err)
	}
	resp, err := c.Patch(neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseMetaport(resp)
}

func DeleteMetaport(c *Client, neID string) (*Metaport, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, metaportEndpoint, neID)
	resp, err := c.Delete(url, nil)
	if err != nil {
		return nil, err
	}
	return parseMetaport(resp)
}
