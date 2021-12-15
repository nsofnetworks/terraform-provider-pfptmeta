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

const enterpriseDNSEndpoint = "v1/enterprise_dns"

type EnterpriseDNS struct {
	ID            string         `json:"id,omitempty"`
	Name          string         `json:"name,omitempty"`
	Description   string         `json:"description,omitempty"`
	MappedDomains []MappedDomain `json:"mapped_domains,omitempty"`
}

func NewEnterpriseDNS(d *schema.ResourceData) *EnterpriseDNS {
	res := &EnterpriseDNS{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	if d.HasChange("mapped_domains") {
		res.MappedDomains = newMappedDomains(d)
	}
	return res
}

func newMappedDomains(d *schema.ResourceData) []MappedDomain {
	mds := d.Get("mapped_domains").([]interface{})
	resp := make([]MappedDomain, len(mds))
	for i, v := range mds {
		md := v.(map[string]interface{})
		resp[i] = MappedDomain{Name: md["name"].(string), MappedDomain: md["mapped_domain"].(string)}
	}
	return resp
}

func parseEnterpriseDNS(resp *http.Response) (*EnterpriseDNS, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	ed := &EnterpriseDNS{}
	err = json.Unmarshal(body, ed)
	if err != nil {
		return nil, fmt.Errorf("could not parse enterprise dns response: %v", err)
	}
	return ed, nil
}

func CreateEnterpriseDNS(ctx context.Context, c *Client, ed *EnterpriseDNS) (*EnterpriseDNS, error) {
	edUrl := fmt.Sprintf("%s/%s", c.BaseURL, enterpriseDNSEndpoint)
	body, err := json.Marshal(ed)
	if err != nil {
		return nil, fmt.Errorf("could not convert enterprise dns to json: %v", err)
	}
	resp, err := c.Post(ctx, edUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseEnterpriseDNS(resp)
}

func UpdateEnterpriseDNS(ctx context.Context, c *Client, edID string, ed *EnterpriseDNS) (*EnterpriseDNS, error) {
	edUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, enterpriseDNSEndpoint, edID)
	body, err := json.Marshal(ed)
	if err != nil {
		return nil, fmt.Errorf("could not convert enterprise dns to json: %v", err)
	}
	resp, err := c.Patch(ctx, edUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseEnterpriseDNS(resp)
}

func GetEnterpriseDNS(ctx context.Context, c *Client, edID string) (*EnterpriseDNS, error) {
	edUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, enterpriseDNSEndpoint, edID)
	resp, err := c.Get(ctx, edUrl, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseEnterpriseDNS(resp)
}

func DeleteEnterpriseDNS(ctx context.Context, c *Client, edID string) (*EnterpriseDNS, error) {
	edUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, enterpriseDNSEndpoint, edID)
	resp, err := c.Delete(ctx, edUrl, nil)
	if err != nil {
		return nil, err
	}
	return parseEnterpriseDNS(resp)
}
