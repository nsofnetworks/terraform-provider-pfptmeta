package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	u "net/url"
)

const cloudAppsEndpoint = "v1/cloud_apps"

type CloudApp struct {
	ID          string   `json:"id,omitempty"`
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description"`
	App         string   `json:"app,omitempty"`
	Tenant      string   `json:"tenant,omitempty"`
	TenantType  string   `json:"tenant_type,omitempty"`
	Type        string   `json:"type,omitempty"`
	Urls        []string `json:"urls"`
}

func NewCloudApp(d *schema.ResourceData) *CloudApp {
	res := &CloudApp{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.App = d.Get("app").(string)
	res.Tenant = d.Get("tenant").(string)
	res.TenantType = d.Get("tenant_type").(string)
	res.Urls = ConfigToStringSlice("urls", d)
	return res
}

func validateCatalogApp(ctx context.Context, c *Client, ca *CloudApp) error {
	catalogApp, err := GetCatalogAppByID(ctx, c, ca.App)
	if err != nil {
		return err
	}
	if ca.Tenant != "" && !catalogApp.Attributes.TenantAwarenessData.TenantCorpIdSupport {
		return fmt.Errorf("catalog app %s (%s) does not support tenant id", catalogApp.Name, catalogApp.ID)
	}
	if ca.TenantType != "All" && !catalogApp.Attributes.TenantAwarenessData.TenantTypeSupport {
		return fmt.Errorf("catalog app %s (%s) does not support tenant type %s", catalogApp.Name, catalogApp.ID, ca.TenantType)
	}
	return nil
}

func parseCloudApp(resp []byte) (*CloudApp, error) {
	c := &CloudApp{}
	err := json.Unmarshal(resp, c)
	if err != nil {
		return nil, fmt.Errorf("could not parse cloud app response: %v", err)
	}
	return c, nil
}

func CreateCloudApp(ctx context.Context, c *Client, ca *CloudApp) (*CloudApp, error) {
	if err := validateCatalogApp(ctx, c, ca); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s", c.BaseURL, cloudAppsEndpoint)
	body, err := json.Marshal(ca)
	if err != nil {
		return nil, fmt.Errorf("could not convert cloud app to json: %v", err)
	}
	resp, err := c.Post(ctx, url, body)
	if err != nil {
		return nil, err
	}
	return parseCloudApp(resp)
}

func UpdateCloudApp(ctx context.Context, c *Client, cID string, ca *CloudApp) (*CloudApp, error) {
	if err := validateCatalogApp(ctx, c, ca); err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, cloudAppsEndpoint, cID)
	body, err := json.Marshal(ca)
	if err != nil {
		return nil, fmt.Errorf("could not convert cloud app to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, body)
	if err != nil {
		return nil, err
	}
	return parseCloudApp(resp)
}

func GetCloudApp(ctx context.Context, c *Client, cID string) (*CloudApp, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, cloudAppsEndpoint, cID)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseCloudApp(resp)
}

func DeleteCloudApp(ctx context.Context, c *Client, cID string) (*CloudApp, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, cloudAppsEndpoint, cID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseCloudApp(resp)
}
