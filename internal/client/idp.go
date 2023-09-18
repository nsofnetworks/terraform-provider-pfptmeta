package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	u "net/url"
)

const idpEndpoint = "v1/settings/idps"

type IdpSamlConfig struct {
	Certificate       string `json:"certificate,omitempty"`
	Issuer            string `json:"issuer,omitempty"`
	SsoUrl            string `json:"sso_url,omitempty"`
	AuthnContextClass string `json:"authn_context_class,omitempty"`
	JitEnabled        bool   `json:"jit_enabled,omitempty"`
}

type IdpOidcConfig struct {
	Issuer       string `json:"issuer,omitempty"`
	ClientId     string `json:"client_id,omitempty"`
	ClientSecret string `json:"client_secret"`
	JitEnabled   bool   `json:"jit_enabled,omitempty"`
}

type IdpScimConfig struct {
	ApiKeyId        string `json:"api_key_id,omitempty"`
	AssumeOwnership bool   `json:"assume_ownership,omitempty"`
}

type Idp struct {
	ID               string         `json:"id,omitempty"`
	Name             string         `json:"name,omitempty"`
	Description      *string        `json:"description"`
	Enabled          bool           `json:"enabled"`
	Hidden           bool           `json:"hidden"`
	Icon             *string        `json:"icon"`
	MappedAttributes *[]string      `json:"mapped_attributes"`
	SamlConfig       *IdpSamlConfig `json:"saml_config,omitempty"`
	OidcConfig       *IdpOidcConfig `json:"oidc_config,omitempty"`
	ScimConfig       *IdpScimConfig `json:"scim_config,omitempty"`
}

func NewIdp(d *schema.ResourceData) *Idp {
	res := &Idp{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Enabled = d.Get("enabled").(bool)
	res.Hidden = d.Get("hidden").(bool)
	icon := d.Get("icon").(string)
	if icon != "" {
		res.Icon = &icon
	} else {
		res.Icon = nil
	}
	dsc := d.Get("description").(string)
	if dsc != "" {
		res.Description = &dsc
	} else {
		res.Description = nil
	}
	ma := ResourceTypeSetToStringSlice(d.Get("mapped_attributes").(*schema.Set))
	if len(ma) > 0 {
		res.MappedAttributes = &ma
	} else {
		res.MappedAttributes = nil
	}
	res.SamlConfig = NewIdpSamlConfig(d)
	res.OidcConfig = NewIdpOidcConfig(d)
	res.ScimConfig = NewIdpScimConfig(d)
	return res
}

func NewIdpSamlConfig(d *schema.ResourceData) *IdpSamlConfig {
	res := &IdpSamlConfig{}
	sc, exists := d.GetOk("saml_config")
	if !exists {
		return nil
	}
	saml := sc.([]interface{})
	if len(saml) != 1 {
		return nil
	}
	saml_conf := saml[0].(map[string]interface{})
	if _, ok := saml_conf["issuer"]; ok {
		res.Issuer = saml_conf["issuer"].(string)
	}
	if _, ok := saml_conf["certificate"]; ok {
		res.Certificate = saml_conf["certificate"].(string)
	}
	if _, ok := saml_conf["sso_url"]; ok {
		res.SsoUrl = saml_conf["sso_url"].(string)
	}
	if _, ok := saml_conf["authn_context_class"]; ok {
		res.AuthnContextClass = saml_conf["authn_context_class"].(string)
	}
	if _, ok := saml_conf["jit_enabled"]; ok {
		res.JitEnabled = saml_conf["jit_enabled"].(bool)
	}
	return res
}

func NewIdpOidcConfig(d *schema.ResourceData) *IdpOidcConfig {
	res := &IdpOidcConfig{}
	oc, exists := d.GetOk("oidc_config")
	if !exists {
		return nil
	}
	oidc := oc.([]interface{})
	if len(oidc) != 1 {
		return nil
	}
	oidc_conf := oidc[0].(map[string]interface{})
	if _, ok := oidc_conf["issuer"]; ok {
		res.Issuer = oidc_conf["issuer"].(string)
	}
	if _, ok := oidc_conf["client_id"]; ok {
		res.ClientId = oidc_conf["client_id"].(string)
	}
	if _, ok := oidc_conf["client_secret"]; ok {
		res.ClientSecret = oidc_conf["client_secret"].(string)
	}
	if _, ok := oidc_conf["jit_enabled"]; ok {
		res.JitEnabled = oidc_conf["jit_enabled"].(bool)
	}
	return res
}

func NewIdpScimConfig(d *schema.ResourceData) *IdpScimConfig {
	res := &IdpScimConfig{}
	sc, exists := d.GetOk("scim_config")
	if !exists {
		return nil
	}
	scim := sc.([]interface{})
	if len(scim) != 1 {
		return nil
	}
	scim_conf := scim[0].(map[string]interface{})
	if _, ok := scim_conf["api_key_id"]; ok {
		res.ApiKeyId = scim_conf["api_key_id"].(string)
	}
	if _, ok := scim_conf["assume_ownership"]; ok {
		res.AssumeOwnership = scim_conf["assume_ownership"].(bool)
	}
	return res
}

func parseIdp(resp []byte) (*Idp, error) {
	idp := &Idp{}
	err := json.Unmarshal(resp, idp)
	if err != nil {
		return nil, fmt.Errorf("could not parse idp response: %v", err)
	}
	return idp, nil
}

func CreateIdp(ctx context.Context, c *Client, idp *Idp) (*Idp, error) {
	body, err := json.Marshal(idp)
	if err != nil {
		return nil, fmt.Errorf("could not convert idp to json: %v", err)
	}
	idpUrl := fmt.Sprintf("%s/%s", c.BaseURL, idpEndpoint)
	resp, err := c.Post(ctx, idpUrl, body)
	if err != nil {
		return nil, err
	}
	return parseIdp(resp)
}

func UpdateIdp(ctx context.Context, c *Client, idpID string, idp *Idp) (*Idp, error) {
	body, err := json.Marshal(idp)
	if err != nil {
		return nil, fmt.Errorf("could not convert idp to json: %v", err)
	}
	idpUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, idpEndpoint, idpID)
	resp, err := c.Patch(ctx, idpUrl, body)
	if err != nil {
		return nil, err
	}
	return parseIdp(resp)
}

func GetIdp(ctx context.Context, c *Client, idpID string) (*Idp, error) {
	idpUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, idpEndpoint, idpID)
	resp, err := c.Get(ctx, idpUrl, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseIdp(resp)
}

func DeleteIdp(ctx context.Context, c *Client, idpID string) (*Idp, error) {
	idpUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, idpEndpoint, idpID)
	resp, err := c.Delete(ctx, idpUrl, nil)
	if err != nil {
		return nil, err
	}
	return parseIdp(resp)
}
