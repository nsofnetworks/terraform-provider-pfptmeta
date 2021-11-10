package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"net/http"
	u "net/url"
)

const (
	networkElementsEndpoint string = "/v1/network_elements"
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

type MappedDomain struct {
	MappedDomain string `json:"mapped_domain"`
	Name         string `json:"name,omitempty"`
}

// ReqBody returns a body with mapped_domain only because the name of the mapped domain should be in the path params only
func (md *MappedDomain) ReqBody() ([]byte, error) {
	md.Name = ""
	return json.Marshal(md)
}

type NetworkElementBody struct {
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description,omitempty"`
	Enabled       *bool    `json:"enabled,omitempty"`
	MappedSubnets []string `json:"mapped_subnets,omitempty"`
	MappedService string   `json:"mapped_service,omitempty"`
	Platform      string   `json:"platform,omitempty"`
	OwnerID       string   `json:"owner_id,omitempty"`
}

func NewNetworkElementBody(d *schema.ResourceData) *NetworkElementBody {
	res := &NetworkElementBody{}
	nameExists := d.HasChange("name")
	if nameExists {
		name := d.Get("name")
		res.Name = name.(string)
	}
	descExists := d.HasChange("description")
	if descExists {
		description := d.Get("description")
		res.Description = description.(string)
	}
	enabledExists := d.HasChange("enabled")
	if enabledExists {
		enabled := d.Get("enabled").(bool)
		res.Enabled = &enabled
	}
	mappedSubnetsExists := d.HasChange("mapped_subnets")
	if mappedSubnetsExists {
		_, mappedSubnets := d.GetChange("mapped_subnets")
		listMappedSubnets := ResourceTypeSetToStringSlice(mappedSubnets.(*schema.Set))
		res.MappedSubnets = listMappedSubnets
	}
	mappedServiceExists := d.HasChange("mapped_service")
	if mappedServiceExists {
		_, mappedService := d.GetChange("mapped_service")
		res.MappedService = mappedService.(string)
	}
	if d.HasChange("platform") {
		res.Platform = d.Get("platform").(string)
	}
	if d.HasChange("owner_id") {
		res.OwnerID = d.Get("owner_id").(string)
	}
	return res
}

type NetworkElementResponse struct {
	NetworkElementBody
	ID            string         `json:"id"`
	OrgID         string         `json:"org_id"`
	DnsName       string         `json:"dns_name"`
	NetID         int            `json:"net_id"`
	AutoAliases   []string       `json:"auto_aliases"`
	Groups        []string       `json:"groups"`
	Type          string         `json:"type"`
	Tags          []Tag          `json:"tags,omitempty"`
	Aliases       []string       `json:"aliases"`
	MappedDomains []MappedDomain `json:"mapped_domains"`
	MappedHosts   []MappedHost   `json:"mapped_hosts"`
}

func parseNetworkElement(resp *http.Response) (*NetworkElementResponse, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	networkElement := &NetworkElementResponse{}
	err = json.Unmarshal(body, networkElement)
	if err != nil {
		return nil, fmt.Errorf("could not parse network element response: %v", err)
	}
	return networkElement, nil
}

func CreateNetworkElement(c *Client, ne *NetworkElementBody) (*NetworkElementResponse, error) {
	neUrl := fmt.Sprintf("%s%s", c.BaseURL, networkElementsEndpoint)
	body, err := json.Marshal(ne)
	if err != nil {
		return nil, fmt.Errorf("could not convert network element to json: %v", err)
	}
	resp, err := c.Post(neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseNetworkElement(resp)
}

func UpdateNetworkElement(c *Client, neId string, ne *NetworkElementBody) (*NetworkElementResponse, error) {
	neUrl := fmt.Sprintf("%s%s/%s", c.BaseURL, networkElementsEndpoint, neId)
	body, err := json.Marshal(ne)
	if err != nil {
		return nil, fmt.Errorf("could not convert network element to json: %v", err)
	}
	resp, err := c.Patch(neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseNetworkElement(resp)
}

func GetNetworkElement(c *Client, neID string) (*NetworkElementResponse, error) {
	url := fmt.Sprintf("%s%s/%s", c.BaseURL, networkElementsEndpoint, neID)
	resp, err := c.Get(url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseNetworkElement(resp)
}

func DeleteNetworkElement(c *Client, neID string) (*NetworkElementResponse, error) {
	url := fmt.Sprintf("%s%s/%s", c.BaseURL, networkElementsEndpoint, neID)
	resp, err := c.Delete(url, nil)
	if err != nil {
		return nil, err
	}
	return parseNetworkElement(resp)
}

func AssignNetworkElementTags(c *Client, neID string, tags []*Tag) error {
	body, err := json.Marshal(tags)
	if err != nil {
		return fmt.Errorf("could not convert network element to json: %v", err)
	}
	url := fmt.Sprintf("%s%s/%s/tags", c.BaseURL, networkElementsEndpoint, neID)
	_, err = c.Put(url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}

func AssignNetworkElementAlias(c *Client, neID, alias string) error {
	url := fmt.Sprintf("%s%s/%s/aliases/%s", c.BaseURL, networkElementsEndpoint, neID, alias)
	_, err := c.Put(url, nil)
	if err != nil {
		return err
	}
	return nil
}

func DeleteNetworkElementAlias(c *Client, neID, alias string) error {
	url := fmt.Sprintf("%s%s/%s/aliases/%s", c.BaseURL, networkElementsEndpoint, neID, alias)
	_, err := c.Delete(url, nil)
	if err != nil {
		return err
	}
	return nil
}

func SetMappedDomain(c *Client, neID string, mappedDomain *MappedDomain) error {
	url := fmt.Sprintf("%s%s/%s/mapped_domains/%s", c.BaseURL, networkElementsEndpoint, neID, mappedDomain.Name)
	body, err := mappedDomain.ReqBody()
	if err != nil {
		return fmt.Errorf("could not convert MappedDomain to json")
	}
	_, err = c.Put(url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}

func DeleteMappedDomain(c *Client, neID, name string) error {
	url := fmt.Sprintf("%s%s/%s/mapped_domains/%s", c.BaseURL, networkElementsEndpoint, neID, name)
	_, err := c.Delete(url, nil)
	if err != nil {
		return err
	}
	return nil
}

func SetMappedHost(c *Client, neID string, mappedHost *MappedHost) error {
	url := fmt.Sprintf("%s%s/%s/mapped_hosts/%s", c.BaseURL, networkElementsEndpoint, neID, mappedHost.Name)
	body, err := mappedHost.ReqBody()
	if err != nil {
		return fmt.Errorf("could not convert MappedHost to json")
	}
	_, err = c.Put(url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}

func DeleteMappedHost(c *Client, neID, name string) error {
	url := fmt.Sprintf("%s%s/%s/mapped_hosts/%s", c.BaseURL, networkElementsEndpoint, neID, name)
	_, err := c.Delete(url, nil)
	if err != nil {
		return err
	}
	return nil
}
