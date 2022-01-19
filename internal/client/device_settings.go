package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"net/http"
)

const deviceSettingsEndpoint = "v1/settings/device"

type DeviceSettings struct {
	ID                        string   `json:"id,omitempty"`
	Name                      string   `json:"name,omitempty"`
	Description               string   `json:"description"`
	Enabled                   bool     `json:"enabled"`
	ApplyOnOrg                bool     `json:"apply_on_org"`
	ApplyToEntities           []string `json:"apply_to_entities"`
	AutoFqdnDomainNames       []string `json:"auto_fqdn_domain_names,omitempty"`
	DirectSso                 *string  `json:"direct_sso,omitempty"`
	OverlayMfaRefreshPeriod   *int     `json:"overlay_mfa_refresh_period,omitempty"`
	OverlayMfaRequired        *bool    `json:"overlay_mfa_required,omitempty"`
	ProtocolSelectionLifetime *int     `json:"protocol_selection_lifetime,omitempty"`
	ProxyAlwaysOn             *bool    `json:"proxy_always_on,omitempty"`
	SearchDomains             []string `json:"search_domains,omitempty"`
	SessionLifetime           *int     `json:"session_lifetime,omitempty"`
	SessionLifetimeGrace      *int     `json:"session_lifetime_grace,omitempty"`
	TunnelMode                *string  `json:"tunnel_mode,omitempty"`
	VpnLoginBrowser           *string  `json:"vpn_login_browser,omitempty"`
	ZtnaAlwaysOn              *bool    `json:"ztna_always_on,omitempty"`
}

func NewDeviceSettings(d *schema.ResourceData) *DeviceSettings {
	res := &DeviceSettings{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Enabled = d.Get("enabled").(bool)
	res.ApplyOnOrg = d.Get("apply_on_org").(bool)

	appliedEntities := d.Get("apply_to_entities").([]interface{})
	res.ApplyToEntities = make([]string, len(appliedEntities))
	for i, val := range appliedEntities {
		res.ApplyToEntities[i] = val.(string)
	}

	domainNames := d.Get("auto_fqdn_domain_names").([]interface{})
	if len(domainNames) == 0 {
		res.AutoFqdnDomainNames = make([]string, len(domainNames))
		for i, val := range domainNames {
			if val != nil {
				res.AutoFqdnDomainNames[i] = val.(string)
			} else {
				res.AutoFqdnDomainNames[i] = ""
			}
		}
	}

	ds, exists := d.GetOk("direct_sso")
	if exists {
		directSso := ds.(string)
		res.DirectSso = &directSso
	}

	omrp, exists := d.GetOk("overlay_mfa_refresh_period")
	if exists {
		overlayMfaRefreshPeriod := omrp.(int)
		res.OverlayMfaRefreshPeriod = &overlayMfaRefreshPeriod
	}

	omr, exists := d.GetOk("overlay_mfa_required")
	if exists {
		overlayMfaRequired := omr.(bool)
		res.OverlayMfaRequired = &overlayMfaRequired
	}

	psl, exists := d.GetOk("protocol_selection_lifetime")
	if exists {
		protocolSelectionLifetime := psl.(int)
		res.ProtocolSelectionLifetime = &protocolSelectionLifetime
	}

	pao, exists := d.GetOk("proxy_always_on")
	if exists {
		proxyAlwaysOn := pao.(bool)
		res.ProxyAlwaysOn = &proxyAlwaysOn
	}

	searchDomains := d.Get("search_domains").([]interface{})
	if len(searchDomains) != 0 {
		res.SearchDomains = make([]string, len(searchDomains))
		for i, val := range searchDomains {
			res.SearchDomains[i] = val.(string)
		}
	}
	slt, exists := d.GetOk("session_lifetime")
	if exists {
		sessionLifetime := slt.(int)
		res.SessionLifetime = &sessionLifetime
	}

	slg, exists := d.GetOk("session_lifetime_grace")
	if exists {
		sessionLifetimeGrace := slg.(int)
		res.SessionLifetimeGrace = &sessionLifetimeGrace
	}

	tm, exists := d.GetOk("tunnel_mode")
	if exists {
		tunnelMode := tm.(string)
		res.TunnelMode = &tunnelMode
	}
	vlb, exists := d.GetOk("vpn_login_browser")
	if exists {
		vpnLoginBrowser := vlb.(string)
		res.VpnLoginBrowser = &vpnLoginBrowser
	}
	ztnaAlwaysOn, exists := d.GetOk("ztna_always_on")
	if exists {
		alwaysOn := ztnaAlwaysOn.(bool)
		res.ZtnaAlwaysOn = &alwaysOn
	} else {
		res.ZtnaAlwaysOn = nil
	}

	return res
}

func parseDeviceSettings(resp *http.Response) (*DeviceSettings, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read device settings response: %v", err)
	}
	a := &DeviceSettings{}
	err = json.Unmarshal(body, a)
	if err != nil {
		return nil, fmt.Errorf("could not parse device settings response: %v", err)
	}
	return a, nil
}

func CreateDeviceSettings(ctx context.Context, c *Client, a *DeviceSettings) (*DeviceSettings, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, deviceSettingsEndpoint)
	body, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("could not convert device settings to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseDeviceSettings(resp)
}

func UpdateDeviceSettings(ctx context.Context, c *Client, aID string, a *DeviceSettings) (*DeviceSettings, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, deviceSettingsEndpoint, aID)
	body, err := json.Marshal(a)
	if err != nil {
		return nil, fmt.Errorf("could not convert device settings to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseDeviceSettings(resp)
}

func GetDeviceSettings(ctx context.Context, c *Client, aID string) (*DeviceSettings, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, deviceSettingsEndpoint, aID)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseDeviceSettings(resp)
}

func DeleteDeviceSettings(ctx context.Context, c *Client, aID string) (*DeviceSettings, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, deviceSettingsEndpoint, aID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseDeviceSettings(resp)
}
