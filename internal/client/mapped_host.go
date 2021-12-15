package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type MappedHost struct {
	MappedHost string `json:"mapped_host"`
	Name       string `json:"name,omitempty"`
}

// ReqBody returns a body with mapped_host only because the name of the mapped host should be in the path params only
func (mh *MappedHost) ReqBody() ([]byte, error) {
	mh.Name = ""
	return json.Marshal(mh)
}

func parseMappedHost(resp *http.Response) (*MappedHost, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	mh := &MappedHost{}
	err = json.Unmarshal(body, mh)
	if err != nil {
		return nil, fmt.Errorf("could not parse network element response: %v", err)
	}
	return mh, nil
}

func SetMappedHost(ctx context.Context, c *Client, neID string, mappedHost *MappedHost) (*MappedHost, error) {
	url := fmt.Sprintf("%s/%s/%s/mapped_hosts/%s", c.BaseURL, networkElementsEndpoint, neID, mappedHost.Name)
	body, err := mappedHost.ReqBody()
	if err != nil {
		return nil, fmt.Errorf("could not convert MappedHost to json")
	}
	resp, err := c.Put(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseMappedHost(resp)
}

func GetMappedHost(ctx context.Context, c *Client, neID string, mappedHost *MappedHost) (*MappedHost, error) {
	url := fmt.Sprintf("%s/%s/%s/mapped_hosts/%s", c.BaseURL, networkElementsEndpoint, neID, mappedHost.Name)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseMappedHost(resp)
}

func DeleteMappedHost(ctx context.Context, c *Client, neID, name string) error {
	url := fmt.Sprintf("%s/%s/%s/mapped_hosts/%s", c.BaseURL, networkElementsEndpoint, neID, name)
	_, err := c.Delete(ctx, url, nil)
	if err != nil {
		return err
	}
	return nil
}
