package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MappedDomain struct {
	MappedDomain string `json:"mapped_domain"`
	Name         string `json:"name,omitempty"`
}

// ReqBody returns a body with mapped_domain only because the name of the mapped domain should be in the path params only
func (md *MappedDomain) ReqBody() ([]byte, error) {
	md.Name = ""
	return json.Marshal(md)
}

func parseMappedDomain(resp *http.Response) (*MappedDomain, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	md := &MappedDomain{}
	err = json.Unmarshal(body, md)
	if err != nil {
		return nil, fmt.Errorf("could not parse network element response: %v", err)
	}
	return md, nil
}

func GetMappedDomain(ctx context.Context, c *Client, neID string, mappedDomain *MappedDomain) (*MappedDomain, error) {
	url := fmt.Sprintf("%s/%s/%s/mapped_domains/%s", c.BaseURL, networkElementsEndpoint, neID, mappedDomain.Name)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseMappedDomain(resp)
}

func SetMappedDomain(ctx context.Context, c *Client, neID string, mappedDomain *MappedDomain) (*MappedDomain, error) {
	url := fmt.Sprintf("%s/%s/%s/mapped_domains/%s", c.BaseURL, networkElementsEndpoint, neID, mappedDomain.Name)
	body, err := mappedDomain.ReqBody()
	if err != nil {
		return nil, fmt.Errorf("could not convert MappedDomain to json")
	}
	resp, err := c.Put(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseMappedDomain(resp)
}

func DeleteMappedDomain(ctx context.Context, c *Client, neID, name string) error {
	url := fmt.Sprintf("%s/%s/%s/mapped_domains/%s", c.BaseURL, networkElementsEndpoint, neID, name)
	_, err := c.Delete(ctx, url, nil)
	if err != nil {
		return err
	}
	return nil
}
