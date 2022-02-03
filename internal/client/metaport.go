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
	metaportEndpoint string = "v1/metaports"
)

type Metaport struct {
	ID                   string   `json:"id,omitempty"`
	Name                 string   `json:"name,omitempty"`
	Description          string   `json:"description,omitempty"`
	Enabled              *bool    `json:"enabled,omitempty"`
	AllowSupport         *bool    `json:"allow_support,omitempty"`
	MappedElements       []string `json:"mapped_elements"`
	NotificationChannels []string `json:"notification_channels"`
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

	res.NotificationChannels = ConfigToStringSlice("notification_channels", d)

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

func CreateMetaport(ctx context.Context, c *Client, m *Metaport) (*Metaport, error) {
	neUrl := fmt.Sprintf("%s/%s", c.BaseURL, metaportEndpoint)
	body, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("could not convert metaport to json: %v", err)
	}
	resp, err := c.Post(ctx, neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseMetaport(resp)
}

func GetMetaport(ctx context.Context, c *Client, mId string) (*Metaport, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, metaportEndpoint, mId)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseMetaport(resp)
}

func GetMetaportByName(ctx context.Context, c *Client, name string) (*Metaport, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, metaportEndpoint)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read metaport response")
	}
	var respBody []Metaport
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return nil, fmt.Errorf("could not parse metaport response: %v", err)
	}
	var nameMatch []Metaport
	for _, m := range respBody {
		if m.Name == name {
			nameMatch = append(nameMatch, m)
		}
	}
	switch len(nameMatch) {
	case 0:
		return nil, fmt.Errorf("could not find metaport with name \"%s\"", name)
	case 1:
		return &nameMatch[0], nil
	default:
		return nil, fmt.Errorf("found more then one metaport with name \"%s\"", name)
	}
}

func UpdateMetaport(ctx context.Context, c *Client, mId string, m *Metaport) (*Metaport, error) {
	neUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, metaportEndpoint, mId)
	body, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("could not convert metaport to json: %v", err)
	}
	resp, err := c.Patch(ctx, neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseMetaport(resp)
}

func DeleteMetaport(ctx context.Context, c *Client, mID string) (*Metaport, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, metaportEndpoint, mID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseMetaport(resp)
}

func AddMappedElementsToMetaport(ctx context.Context, c *Client, mID string, meIDs []string) (*Metaport, error) {
	url := fmt.Sprintf("%s/%s/%s/add_mapped_elements", c.BaseURL, metaportEndpoint, mID)
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
	return parseMetaport(resp)
}

func RemoveMappedElementsFromMetaport(ctx context.Context, c *Client, mID string, meIDs []string) (*Metaport, error) {
	url := fmt.Sprintf("%s/%s/%s/remove_mapped_elements", c.BaseURL, metaportEndpoint, mID)
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
	return parseMetaport(resp)
}
