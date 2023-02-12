package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	tenantRestrictionEndpoint string = "v1/tenant_restrictions"
)

type GoogleConfig struct {
	AllowConsumerAccess  bool     `json:"allow_consumer_access"`
	AllowServiceAccounts bool     `json:"allow_service_accounts"`
	Tenants              []string `json:"tenants"`
}

func newGoogleConfig(d *schema.ResourceData) *GoogleConfig {
	res := &GoogleConfig{}
	rawGoogleConfig := d.Get("google_config").([]interface{})
	if len(rawGoogleConfig) == 0 {
		return nil
	}
	googleConfig := rawGoogleConfig[0].(map[string]interface{})
	tenants := googleConfig["tenants"].([]interface{})
	res.Tenants = make([]string, len(tenants))
	for i, val := range tenants {
		res.Tenants[i] = val.(string)
	}
	res.AllowServiceAccounts = googleConfig["allow_service_accounts"].(bool)
	res.AllowConsumerAccess = googleConfig["allow_consumer_access"].(bool)
	return res
}

type MicrosoftConfig struct {
	AllowPersonalMicrosoftDomains bool     `json:"allow_personal_microsoft_domains"`
	TenantDirectoryId             string   `json:"tenant_directory_id"`
	Tenants                       []string `json:"tenants"`
}

func newMicrosoftConfig(d *schema.ResourceData) *MicrosoftConfig {
	res := &MicrosoftConfig{}
	rawMicrosoftConfig := d.Get("microsoft_config").([]interface{})
	if len(rawMicrosoftConfig) == 0 {
		return nil
	}
	googleConfig := rawMicrosoftConfig[0].(map[string]interface{})
	tenants := googleConfig["tenants"].([]interface{})
	res.Tenants = make([]string, len(tenants))
	for i, val := range tenants {
		res.Tenants[i] = val.(string)
	}
	res.AllowPersonalMicrosoftDomains = googleConfig["allow_personal_microsoft_domains"].(bool)
	res.TenantDirectoryId = googleConfig["tenant_directory_id"].(string)
	return res
}

type TenantRestriction struct {
	ID              string           `json:"id,omitempty"`
	Name            string           `json:"name,omitempty"`
	Description     string           `json:"description"`
	GoogleConfig    *GoogleConfig    `json:"google_config,omitempty"`
	MicrosoftConfig *MicrosoftConfig `json:"microsoft_config,omitempty"`
	Type            string           `json:"type,omitempty"`
}

func NewTenantRestriction(d *schema.ResourceData) *TenantRestriction {
	res := &TenantRestriction{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.GoogleConfig = newGoogleConfig(d)
	res.MicrosoftConfig = newMicrosoftConfig(d)
	return res
}

func parseTenantRestriction(resp []byte) (*TenantRestriction, error) {
	nc := &TenantRestriction{}
	err := json.Unmarshal(resp, nc)
	if err != nil {
		return nil, fmt.Errorf("could not parse tenant restriction response: %v", err)
	}
	return nc, nil
}

func CreateTenantRestriction(ctx context.Context, c *Client, tr *TenantRestriction) (*TenantRestriction, error) {
	trUrl := fmt.Sprintf("%s/%s", c.BaseURL, tenantRestrictionEndpoint)
	body, err := json.Marshal(tr)
	if err != nil {
		return nil, fmt.Errorf("could not convert tenant restriction to json: %v", err)
	}
	resp, err := c.Post(ctx, trUrl, body)
	if err != nil {
		return nil, err
	}
	return parseTenantRestriction(resp)
}

func UpdateTenantRestriction(ctx context.Context, c *Client, trID string, tr *TenantRestriction) (*TenantRestriction, error) {
	ncUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, tenantRestrictionEndpoint, trID)
	body, err := json.Marshal(tr)
	if err != nil {
		return nil, fmt.Errorf("could not convert tenant restriction to json: %v", err)
	}
	resp, err := c.Patch(ctx, ncUrl, body)
	if err != nil {
		return nil, err
	}
	return parseTenantRestriction(resp)
}

func GetTenantRestriction(ctx context.Context, c *Client, trID string) (*TenantRestriction, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, tenantRestrictionEndpoint, trID)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseTenantRestriction(resp)
}

func DeleteTenantRestriction(ctx context.Context, c *Client, ncID string) (*TenantRestriction, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, tenantRestrictionEndpoint, ncID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseTenantRestriction(resp)
}
