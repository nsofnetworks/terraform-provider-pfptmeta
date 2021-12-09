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

const groupEndpoint string = "v1/groups"

type Group struct {
	ID            string   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description"`
	Expression    *string  `json:"expression"`
	ProvisionedBy string   `json:"provisioned_by,omitempty"`
	Roles         []string `json:"roles,omitempty"`
}

func NewGroup(d *schema.ResourceData) *Group {
	res := &Group{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	e := d.Get("expression").(string)
	if e == "" {
		res.Expression = nil
	} else {
		res.Expression = &e
	}
	return res
}

func parseGroup(resp *http.Response) (*Group, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	g := &Group{}
	err = json.Unmarshal(body, g)
	if err != nil {
		return nil, fmt.Errorf("could not parse group response: %v", err)
	}
	return g, nil
}

func CreateGroup(c *Client, g *Group) (*Group, error) {
	gUrl := fmt.Sprintf("%s/%s", c.BaseURL, groupEndpoint)
	body, err := json.Marshal(g)
	if err != nil {
		return nil, fmt.Errorf("could not convert group to json: %v", err)
	}
	resp, err := c.Post(gUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseGroup(resp)
}

func UpdateGroup(c *Client, gID string, g *Group) (*Group, error) {
	neUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, groupEndpoint, gID)
	body, err := json.Marshal(g)
	if err != nil {
		return nil, fmt.Errorf("could not convert group to json: %v", err)
	}
	resp, err := c.Patch(neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseGroup(resp)
}

func GetGroupById(c *Client, gID string) (*Group, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, groupEndpoint, gID)
	resp, err := c.Get(url, nil)
	if err != nil {
		return nil, err
	}
	return parseGroup(resp)
}
func GetGroupByName(c *Client, name string) (*Group, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, groupEndpoint)
	resp, err := c.Get(url, u.Values{"name": {name}, "pagination": {"false"}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	gs := &[]Group{}
	err = json.Unmarshal(body, gs)
	if err != nil {
		return nil, fmt.Errorf("could not parse group response: %v", err)
	}
	for _, g := range *gs {
		if g.Name == name {
			return &g, nil
		}
	}
	return nil, nil
}

func DeleteGroup(c *Client, gID string) (*Group, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, groupEndpoint, gID)
	resp, err := c.Delete(url, nil)
	if err != nil {
		return nil, err
	}
	return parseGroup(resp)
}

func AssignRolesToGroup(c *Client, gID string, roles []string) ([]string, error) {
	url := fmt.Sprintf("%s/%s/%s/roles", c.BaseURL, groupEndpoint, gID)
	body, err := json.Marshal(roles)
	if err != nil {
		return nil, fmt.Errorf("could not convert roles to json: %v", err)
	}
	resp, err := c.Put(url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	g, err := parseGroup(resp)
	if err != nil {
		return nil, err
	}
	return g.Roles, nil
}
