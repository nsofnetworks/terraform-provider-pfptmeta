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

const groupEndpoint string = "v1/groups"

type Group struct {
	ID            string   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description"`
	Expression    *string  `json:"expression"`
	ProvisionedBy string   `json:"provisioned_by,omitempty"`
	Roles         []string `json:"roles,omitempty"`
	Users         []string `json:"users,omitempty"`
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

func CreateGroup(ctx context.Context, c *Client, g *Group) (*Group, error) {
	gUrl := fmt.Sprintf("%s/%s", c.BaseURL, groupEndpoint)
	body, err := json.Marshal(g)
	if err != nil {
		return nil, fmt.Errorf("could not convert group to json: %v", err)
	}
	resp, err := c.Post(ctx, gUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseGroup(resp)
}

func UpdateGroup(ctx context.Context, c *Client, gID string, g *Group) (*Group, error) {
	neUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, groupEndpoint, gID)
	body, err := json.Marshal(g)
	if err != nil {
		return nil, fmt.Errorf("could not convert group to json: %v", err)
	}
	resp, err := c.Patch(ctx, neUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseGroup(resp)
}

func GetGroupById(ctx context.Context, c *Client, gID string) (*Group, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, groupEndpoint, gID)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseGroup(resp)
}
func GetGroupByName(ctx context.Context, c *Client, name string) (*Group, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, groupEndpoint)
	resp, err := c.Get(ctx, url, u.Values{"name": {name}, "pagination": {"false"}})
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

func DeleteGroup(ctx context.Context, c *Client, gID string) (*Group, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, groupEndpoint, gID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseGroup(resp)
}

func AssignRolesToGroup(ctx context.Context, c *Client, gID string, roles []string) ([]string, error) {
	url := fmt.Sprintf("%s/%s/%s/roles", c.BaseURL, groupEndpoint, gID)
	body, err := json.Marshal(roles)
	if err != nil {
		return nil, fmt.Errorf("could not convert roles to json: %v", err)
	}
	resp, err := c.Put(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	g, err := parseGroup(resp)
	if err != nil {
		return nil, err
	}
	return g.Roles, nil
}

func AddUsersToGroup(ctx context.Context, c *Client, gID string, users []string) error {
	url := fmt.Sprintf("%s/%s/%s/add_users", c.BaseURL, groupEndpoint, gID)
	body, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("could not convert users to json: %v", err)
	}
	_, err = c.Post(ctx, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}

func RemoveUsersFromGroup(ctx context.Context, c *Client, gID string, users []string) error {
	url := fmt.Sprintf("%s/%s/%s/remove_users", c.BaseURL, groupEndpoint, gID)
	body, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("could not convert users to json: %v", err)
	}
	_, err = c.Post(ctx, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}
