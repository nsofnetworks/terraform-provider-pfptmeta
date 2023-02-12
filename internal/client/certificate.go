package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	u "net/url"
)

const certificateEndpoint = "v1/certificates"

type Certificate struct {
	ID                string   `json:"id,omitempty"`
	Name              string   `json:"name,omitempty"`
	Description       string   `json:"description"`
	Sans              []string `json:"sans,omitempty"`
	SerialNumber      string   `json:"serial_number,omitempty"`
	Status            string   `json:"status,omitempty"`
	StatusDescription string   `json:"status_description,omitempty"`
	ValidNotAfter     string   `json:"valid_not_after,omitempty"`
	ValidNotBefore    string   `json:"valid_not_before,omitempty"`
}

func NewCertificate(d *schema.ResourceData) *Certificate {
	res := &Certificate{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	if d.HasChange("sans") {
		res.Sans = ResourceTypeSetToStringSlice(d.Get("sans").(*schema.Set))
	}
	return res
}

func parseCertificate(resp []byte) (*Certificate, error) {
	c := &Certificate{}
	err := json.Unmarshal(resp, c)
	if err != nil {
		return nil, fmt.Errorf("could not parse certificate response: %v", err)
	}
	return c, nil
}

func CreateCertificate(ctx context.Context, c *Client, cert *Certificate) (*Certificate, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, certificateEndpoint)
	body, err := json.Marshal(cert)
	if err != nil {
		return nil, fmt.Errorf("could not convert certificate to json: %v", err)
	}
	resp, err := c.Post(ctx, url, body)
	if err != nil {
		return nil, err
	}
	return parseCertificate(resp)
}

func UpdateCertificate(ctx context.Context, c *Client, cID string, cert *Certificate) (*Certificate, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, certificateEndpoint, cID)
	body, err := json.Marshal(cert)
	if err != nil {
		return nil, fmt.Errorf("could not convert certificate to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, body)
	if err != nil {
		return nil, err
	}
	return parseCertificate(resp)
}

func GetCertificate(ctx context.Context, c *Client, cID string) (*Certificate, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, certificateEndpoint, cID)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseCertificate(resp)
}

func DeleteCertificate(ctx context.Context, c *Client, cID string) (*Certificate, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, certificateEndpoint, cID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseCertificate(resp)
}
