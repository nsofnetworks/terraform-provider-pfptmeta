package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	u "net/url"
)

const easyLinkEndpoint = "v1/easylinks"

type Proxy struct {
	EnterpriseAccess    bool     `json:"enterprise_access"`
	Hosts               []string `json:"hosts"`
	HttpHostHeader      string   `json:"http_host_header"`
	RewriteContentTypes []string `json:"rewrite_content_types"`
	RewriteHosts        bool     `json:"rewrite_hosts"`
	RewriteHostsClient  bool     `json:"rewrite_hosts_client"`
	RewriteHttp         bool     `json:"rewrite_http"`
	SharedCookies       bool     `json:"shared_cookies"`
}

func NewProxy(d *schema.ResourceData) *Proxy {
	res := &Proxy{}
	p, exists := d.GetOk("proxy")
	if !exists || len(p.([]interface{})) != 1 {
		return nil
	}
	proxy := p.([]interface{})[0].(map[string]interface{})
	res.EnterpriseAccess = proxy["enterprise_access"].(bool)
	hosts := proxy["hosts"].([]interface{})
	res.Hosts = make([]string, len(hosts))
	for i, val := range hosts {
		res.Hosts[i] = val.(string)
	}
	res.HttpHostHeader = proxy["http_host_header"].(string)
	rewriteContentTypes := proxy["rewrite_content_types"].([]interface{})
	res.RewriteContentTypes = make([]string, len(rewriteContentTypes))
	for i, val := range rewriteContentTypes {
		res.RewriteContentTypes[i] = val.(string)
	}
	res.RewriteHosts = proxy["rewrite_hosts"].(bool)
	res.RewriteHostsClient = proxy["rewrite_hosts_client"].(bool)
	res.RewriteHttp = proxy["rewrite_http"].(bool)
	res.SharedCookies = proxy["shared_cookies"].(bool)
	return res
}

type Rdp struct {
	RemoteApp            string `json:"remote_app"`
	RemoteAppCmdArgs     string `json:"remote_app_cmd_args"`
	RemoteAppWorkDir     string `json:"remote_app_work_dir"`
	Security             string `json:"security"`
	ServerKeyboardLayout string `json:"server_keyboard_layout"`
}

func NewRdp(d *schema.ResourceData) *Rdp {
	res := &Rdp{}
	r, exists := d.GetOk("rdp")
	if !exists || len(r.([]interface{})) != 1 {
		return nil
	}
	rdp := r.([]interface{})[0].(map[string]interface{})
	res.RemoteApp = rdp["remote_app"].(string)
	res.RemoteAppCmdArgs = rdp["remote_app_cmd_args"].(string)
	res.RemoteAppWorkDir = rdp["remote_app_work_dir"].(string)
	res.Security = rdp["security"].(string)
	res.ServerKeyboardLayout = rdp["server_keyboard_layout"].(string)
	return res
}

type EasyLink struct {
	ID              string   `json:"id,omitempty"`
	Name            string   `json:"name,omitempty"`
	Description     string   `json:"description"`
	DomainName      string   `json:"domain_name,omitempty"`
	Viewers         []string `json:"viewers,omitempty"`
	AccessFqdn      string   `json:"access_fqdn,omitempty"`
	AccessType      string   `json:"access_type"`
	MappedElementId *string  `json:"mapped_element_id"`
	CertificateId   string   `json:"certificate_id,omitempty"`
	Audit           bool     `json:"audit"`
	EnableSni       bool     `json:"enable_sni"`
	Port            int      `json:"port"`
	Protocol        string   `json:"protocol,omitempty"`
	RootPath        string   `json:"root_path"`
	Proxy           *Proxy   `json:"proxy,omitempty"`
	Rdp             *Rdp     `json:"rdp,omitempty"`
	Version         int      `json:"version,omitempty"`
}

func NewEasyLink(d *schema.ResourceData) *EasyLink {
	res := &EasyLink{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.DomainName = d.Get("domain_name").(string)
	res.Viewers = ResourceTypeSetToStringSlice(d.Get("viewers").(*schema.Set))
	res.AccessFqdn = d.Get("access_fqdn").(string)
	res.AccessType = d.Get("access_type").(string)
	me := d.Get("mapped_element_id").(string)
	if me == "" {
		res.MappedElementId = nil
	} else {
		res.MappedElementId = &me
	}
	res.CertificateId = d.Get("certificate_id").(string)
	res.Audit = d.Get("audit").(bool)
	res.EnableSni = d.Get("enable_sni").(bool)
	res.Port = d.Get("port").(int)
	res.Protocol = d.Get("protocol").(string)
	res.RootPath = d.Get("root_path").(string)

	return res
}

func parseEasyLink(resp []byte) (*EasyLink, error) {
	e := &EasyLink{}
	err := json.Unmarshal(resp, e)
	if err != nil {
		return nil, fmt.Errorf("could not parse easy link response: %v", err)
	}
	return e, nil
}

func CreateEasyLink(ctx context.Context, c *Client, e *EasyLink) (*EasyLink, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, easyLinkEndpoint)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("could not convert easy link to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseEasyLink(resp)
}

func GetEasyLink(ctx context.Context, c *Client, eID string) (*EasyLink, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, easyLinkEndpoint, eID)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseEasyLink(resp)
}

func UpdateEasyLink(ctx context.Context, c *Client, eID string, e *EasyLink) (*EasyLink, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, easyLinkEndpoint, eID)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("could not convert easy link to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseEasyLink(resp)
}

func DeleteEasyLink(ctx context.Context, c *Client, eID string) (*EasyLink, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, easyLinkEndpoint, eID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseEasyLink(resp)
}

func UpdateEasylinkProxy(ctx context.Context, c *Client, eID string, p *Proxy) error {
	url := fmt.Sprintf("%s/%s/%s/proxy", c.BaseURL, easyLinkEndpoint, eID)
	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("could not convert proxy to json: %v", err)
	}
	_, err = c.Patch(ctx, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}

func UpdateEasylinkRdp(ctx context.Context, c *Client, eID string, r *Rdp) error {
	url := fmt.Sprintf("%s/%s/%s/rdp", c.BaseURL, easyLinkEndpoint, eID)
	body, err := json.Marshal(r)
	if err != nil {
		return fmt.Errorf("could not convert rdp to json: %v", err)
	}
	_, err = c.Patch(ctx, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}
