package client

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	tunnelEndpoint string = "v1/tunnels"
)

type GreTunnelConfig struct {
	SourceIps []string `json:"source_ips"`
}

type Tunnel struct {
	ID          string           `json:"id,omitempty"`
	Name        string           `json:"name,omitempty"`
	Description string           `json:"description,omitempty"`
	Enabled     *bool            `json:"enabled,omitempty"`
	GreConfig   *GreTunnelConfig `json:"gre_config,omitempty"`
}

func ep(c *Client, tId *string) string {
	ret := fmt.Sprintf("%s/%s", c.BaseURL, tunnelEndpoint)
	if tId != nil {
		ret = fmt.Sprintf("%s/%s", ret, *tId)
	}
	return ret
}

func parseTunnel(resp []byte) (*Tunnel, error) {
	t := &Tunnel{}
	err := json.Unmarshal(resp, t)
	if err != nil {
		return nil, fmt.Errorf("could not parse tunnel response: %v",
			err)
	}
	if t.GreConfig != nil && len(t.GreConfig.SourceIps) == 0 {
		t.GreConfig = nil
	}
	return t, nil
}

func tunnelJsonMarshal(t *Tunnel) ([]byte, error) {
	tcopy := *t
	tcopy.GreConfig = nil // this field is readonly
	b, err := json.Marshal(tcopy)
	if err != nil {
		return nil, fmt.Errorf("could not convert tunnel to json: %v",
			err)
	}
	return b, nil
}

func CreateTunnel(ctx context.Context, c *Client, t *Tunnel) (*Tunnel, error) {
	body, err := tunnelJsonMarshal(t)
	if err != nil {
		return nil, err
	}
	resp, err := c.Post(ctx, ep(c, nil), body)
	if err != nil {
		return nil, err
	}
	return parseTunnel(resp)
}

func GetTunnel(ctx context.Context, c *Client, tId string) (*Tunnel, error) {
	resp, err := c.Get(ctx, ep(c, &tId), nil)
	if err != nil {
		return nil, err
	}
	return parseTunnel(resp)
}

func GetTunnelByName(ctx context.Context, c *Client, name string) (*Tunnel, error) {
	resp, err := c.Get(ctx, ep(c, nil), nil)
	if err != nil {
		return nil, err
	}
	var respBody []Tunnel
	err = json.Unmarshal(resp, &respBody)
	if err != nil {
		return nil, fmt.Errorf("could not parse tunnel response: %v", err)
	}
	var nameMatch []Tunnel
	for _, t := range respBody {
		if t.Name == name {
			nameMatch = append(nameMatch, t)
		}
	}
	switch len(nameMatch) {
	case 0:
		return nil, fmt.Errorf("could not find tunnel with name \"%s\"", name)
	case 1:
		return &nameMatch[0], nil
	default:
		return nil, fmt.Errorf("found more than one tunnel with name \"%s\"", name)
	}
}

func UpdateTunnel(ctx context.Context, c *Client, tId string, t *Tunnel) (*Tunnel, error) {
	body, err := tunnelJsonMarshal(t)
	if err != nil {
		return nil, err
	}
	resp, err := c.Patch(ctx, ep(c, &tId), body)
	if err != nil {
		return nil, err
	}
	return parseTunnel(resp)
}

func DeleteTunnel(ctx context.Context, c *Client, tId string) (*Tunnel, error) {
	resp, err := c.Delete(ctx, ep(c, &tId), nil)
	if err != nil {
		return nil, err
	}
	return parseTunnel(resp)
}

func AddGreSourceIpsToTunnel(ctx context.Context, c *Client, tId string,
	sourceIps []string) (*Tunnel, error) {
	url := fmt.Sprintf("%s/add_gre_source_ip", ep(c, &tId))
	var resp []byte

	body := make(map[string]string)
	for _, ip := range sourceIps {
		body["ip"] = ip
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("could not convert source IP to json: %v", err)
		}
		resp, err = c.Post(ctx, url, jsonBody)
		if err != nil {
			return nil, err
		}
	}
	if resp == nil {
		return nil, fmt.Errorf("No addresses were added")
	}
	return parseTunnel(resp)
}

func RemoveGreSourceIpsFromTunnel(ctx context.Context, c *Client, tId string,
	sourceIps []string) (*Tunnel, error) {
	url := fmt.Sprintf("%s/remove_gre_source_ip", ep(c, &tId))
	var resp []byte

	body := make(map[string]string)
	for _, ip := range sourceIps {
		body["ip"] = ip
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("could not convert source IP to json: %v", err)
		}
		resp, err = c.Post(ctx, url, jsonBody)
		if err != nil {
			return nil, err
		}
	}
	if resp == nil {
		return nil, fmt.Errorf("No addresses were removed")
	}
	return parseTunnel(resp)
}
