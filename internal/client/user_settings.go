package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"strconv"
)

const userSettingsEndpoint = "v1/settings/user"

type UserSettings struct {
	ID                 string   `json:"id,omitempty"`
	Name               string   `json:"name,omitempty"`
	Description        string   `json:"description"`
	Enabled            bool     `json:"enabled"`
	ApplyOnOrg         bool     `json:"apply_on_org"`
	ApplyToEntities    []string `json:"apply_to_entities"`
	AllowedFactors     []string `json:"allowed_factors,omitempty"`
	MaxDevicesPerUser  *int     `json:"max_devices_per_user,omitempty"`
	MfaRequired        *bool    `json:"mfa_required,omitempty"`
	PasswordExpiration *int     `json:"password_expiration,omitempty"`
	ProhibitedOs       []string `json:"prohibited_os,omitempty"`
	ProxyPops          *string  `json:"proxy_pops,omitempty"`
	SsoMandatory       *bool    `json:"sso_mandatory,omitempty"`
}

func NewUserSettings(d *schema.ResourceData) *UserSettings {
	res := &UserSettings{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Enabled = d.Get("enabled").(bool)
	res.ApplyOnOrg = d.Get("apply_on_org").(bool)
	res.ApplyToEntities = ConfigToStringSlice("apply_to_entities", d)
	_, exists := d.GetOk("allowed_factors")
	if exists {
		res.AllowedFactors = ConfigToStringSlice("allowed_factors", d)
	}
	mdpu := d.Get("max_devices_per_user")
	exists = mdpu.(string) != ""
	if exists {
		maxDevicesPerUser, _ := strconv.Atoi(mdpu.(string))
		res.MaxDevicesPerUser = &maxDevicesPerUser
	} else {
		res.MaxDevicesPerUser = nil
	}
	mfar, exists := d.GetOkExists("mfa_required")
	if exists {
		mfaRequired := mfar.(bool)
		res.MfaRequired = &mfaRequired
	}
	pw, exists := d.GetOk("password_expiration")
	if exists {
		passwordExpiration := pw.(int)
		res.PasswordExpiration = &passwordExpiration
	}
	_, exists = d.GetOk("prohibited_os")
	if exists {
		res.ProhibitedOs = ConfigToStringSlice("prohibited_os", d)
	}
	pp, exists := d.GetOk("proxy_pops")
	if exists {
		proxyPops := pp.(string)
		res.ProxyPops = &proxyPops
	}
	ssoM, exists := d.GetOkExists("sso_mandatory")
	if exists {
		ssoMandatory := ssoM.(bool)
		res.SsoMandatory = &ssoMandatory
	}
	return res
}

func parseUserSettings(resp []byte) (*UserSettings, error) {
	ds := &UserSettings{}
	err := json.Unmarshal(resp, ds)
	if err != nil {
		return nil, fmt.Errorf("could not parse user settings response: %v", err)
	}
	return ds, nil
}

func CreateUserSettings(ctx context.Context, c *Client, ds *UserSettings) (*UserSettings, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, userSettingsEndpoint)
	body, err := json.Marshal(ds)
	if err != nil {
		return nil, fmt.Errorf("could not convert user settings to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseUserSettings(resp)
}

func UpdateUserSettings(ctx context.Context, c *Client, dsID string, ds *UserSettings) (*UserSettings, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, userSettingsEndpoint, dsID)
	body, err := json.Marshal(ds)
	if err != nil {
		return nil, fmt.Errorf("could not convert user settings to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseUserSettings(resp)
}

func GetUserSettings(ctx context.Context, c *Client, dsID string) (*UserSettings, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, userSettingsEndpoint, dsID)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseUserSettings(resp)
}

func DeleteUserSettings(ctx context.Context, c *Client, dsID string) (*UserSettings, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, userSettingsEndpoint, dsID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseUserSettings(resp)
}
