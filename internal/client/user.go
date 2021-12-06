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

const UsersEndpoint = "v1/users"

type User struct {
	ID          string  `json:"id,omitempty"`
	GivenName   string  `json:"given_name,omitempty"`
	FamilyName  string  `json:"family_name,omitempty"`
	Description string  `json:"description"`
	Email       string  `json:"email,omitempty"`
	Enabled     *bool   `json:"enabled,omitempty"`
	Phone       *string `json:"phone"`
	Name        string  `json:"name,omitempty"`
	Tags        []Tag   `json:"tags,omitempty"`
}

func NewUser(d *schema.ResourceData) *User {
	res := &User{}
	if d.HasChange("given_name") {
		res.GivenName = d.Get("given_name").(string)
	}
	if d.HasChange("family_name") {
		res.FamilyName = d.Get("family_name").(string)
	}
	res.Description = d.Get("description").(string)

	if d.HasChange("email") {
		res.Email = d.Get("email").(string)
	}
	enabled := d.Get("enabled").(bool)
	res.Enabled = &enabled

	p := d.Get("phone").(string)
	if p == "" {
		res.Phone = nil
	} else {
		res.Phone = &p
	}

	return res
}

func parseUser(resp *http.Response) (*User, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	user := &User{}
	err = json.Unmarshal(body, user)
	if err != nil {
		return nil, fmt.Errorf("could not parse user response: %v", err)
	}
	return user, nil
}

func CreateUser(c *Client, ed *User) (*User, error) {
	uUrl := fmt.Sprintf("%s/%s", c.BaseURL, UsersEndpoint)
	body, err := json.Marshal(ed)
	if err != nil {
		return nil, fmt.Errorf("could not convert user to json: %v", err)
	}
	resp, err := c.Post(uUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseUser(resp)
}

func UpdateUser(c *Client, edID string, ed *User) (*User, error) {
	uUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, UsersEndpoint, edID)
	body, err := json.Marshal(ed)
	if err != nil {
		return nil, fmt.Errorf("could not convert user to json: %v", err)
	}
	resp, err := c.Patch(uUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseUser(resp)
}

func GetUserByID(c *Client, uID string) (*User, error) {
	uUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, UsersEndpoint, uID)
	resp, err := c.Get(uUrl, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseUser(resp)
}

func GetUserByEmail(c *Client, email string) (*User, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, UsersEndpoint)
	resp, err := c.Get(url, u.Values{"expand": {"true"}, "pagination": {"false"}, "email": {email}})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	users := &[]User{}
	err = json.Unmarshal(body, users)
	if err != nil {
		return nil, fmt.Errorf("could not parse user response: %v", err)
	}
	for _, r := range *users {
		if r.Email == email {
			return &r, nil
		}
	}
	return nil, nil
}

func DeleteUser(c *Client, uID string) (*User, error) {
	uUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, UsersEndpoint, uID)
	resp, err := c.Delete(uUrl, nil)
	if err != nil {
		return nil, err
	}
	return parseUser(resp)
}
