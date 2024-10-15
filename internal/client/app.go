package client

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	u "net/url"
)

const appEndpoint = "v1/apps"

type AppSaml struct {
	AudienceUri            string `json:"audience_uri,omitempty"`
	Recipient              string `json:"recipient,omitempty"`
	Destination            string `json:"destination,omitempty"`
	SsoAcsUrl              string `json:"sso_acs_url,omitempty"`
	SubjectNameIdAttribute string `json:"subject_name_id_attribute,omitempty"`
	SubjectNameIdFormat    string `json:"subject_name_id_format,omitempty"`
	SignatureAlgorithm     string `json:"signature_algorithm,omitempty"`
	DigestAlgorithm        string `json:"digest_algorithm,omitempty"`
	DefaultRelayState      string `json:"default_relay_state,omitempty"`
	X509Cert               string `json:"x509_cert,omitempty"`
	IdpIssuer              string `json:"idp_issuer,omitempty"`
	IdpSsoUrl              string `json:"idp_sso_url,omitempty"`
	AuthnContextClass      string `json:"authn_context_class,omitempty"`
}

type AppOidc struct {
	SignInRedirectUrls  []string `json:"sign_in_redirect_urls,omitempty"`
	GrantTypes          []string `json:"grant_types,omitempty"`
	Scopes              []string `json:"scopes,omitempty"`
	InitiateLoginUrl    *string  `json:"initiate_login_url"`
	AccessTokenLifetime int      `json:"access_token_lifetime,omitempty"`
	IdTokenLifetime     int      `json:"id_token_lifetime,omitempty"`
}

type AppMappedAttributes struct {
	AttributeFormat    *string `json:"attribute_format,omitempty"`
	VariableName       string  `json:"variable_name,omitempty"`
	TargetVariableName *string `json:"target_variable_name"`
	FilterType         *string `json:"filter_type"`
	FilterValue        *string `json:"filter_value"`
}

type AppDomainFederation struct {
	Domain string `json:"domain"`
}

type App struct {
	ID               string                `json:"id,omitempty"`
	Name             string                `json:"name,omitempty"`
	Description      *string               `json:"description"`
	Enabled          bool                  `json:"enabled"`
	Visible          bool                  `json:"visible"`
	AssignedMembers  []string              `json:"assigned_members"`
	IpWhitelist      *[]string             `json:"ip_whitelist"`
	DirectSsoLogin   *string               `json:"direct_sso_login"`
	Protocol         string                `json:"protocol,omitempty"`
	Saml             *AppSaml              `json:"saml,omitempty"`
	Oidc             *AppOidc              `json:"oidc,omitempty"`
	MappedAttributes []AppMappedAttributes `json:"mapped_attributes,omitempty"`
	DomainFederation *AppDomainFederation  `json:"domain_federation,omitempty"`
}

func NewApp(d *schema.ResourceData) *App {
	res := &App{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Enabled = d.Get("enabled").(bool)
	res.Visible = d.Get("visible").(bool)
	res.AssignedMembers = ResourceTypeSetToStringSlice(d.Get("assigned_members").(*schema.Set))
	res.Protocol = d.Get("protocol").(string)
	res.Oidc = nil
	res.Saml = nil

	dsl := d.Get("direct_sso_login").(string)
	if dsl != "" {
		res.DirectSsoLogin = &dsl
	} else {
		res.DirectSsoLogin = nil
	}
	dsc := d.Get("description").(string)
	if dsc != "" {
		res.Description = &dsc
	} else {
		res.Description = nil
	}
	ipwl := ResourceTypeSetToStringSlice(d.Get("ip_whitelist").(*schema.Set))
	if len(ipwl) > 0 {
		res.IpWhitelist = &ipwl
	} else {
		res.IpWhitelist = nil
	}
	return res
}

func NewAppSaml(d *schema.ResourceData) *AppSaml {
	res := &AppSaml{}
	s, exists := d.GetOk("saml")
	if !exists {
		return nil
	}
	saml := s.([]interface{})
	if len(saml) != 1 {
		return nil
	}
	saml_conf := saml[0].(map[string]interface{})
	res.AudienceUri = saml_conf["audience_uri"].(string)
	res.Recipient = saml_conf["recipient"].(string)
	res.Destination = saml_conf["destination"].(string)
	res.SsoAcsUrl = saml_conf["sso_acs_url"].(string)
	res.SubjectNameIdAttribute = saml_conf["subject_name_id_attribute"].(string)
	if _, ok := saml_conf["subject_name_id_format"]; ok {
		res.SubjectNameIdFormat = saml_conf["subject_name_id_format"].(string)
	}
	if _, ok := saml_conf["signature_algorithm"]; ok {
		res.SignatureAlgorithm = saml_conf["signature_algorithm"].(string)
	}
	if _, ok := saml_conf["digest_algorithm"]; ok {
		res.DigestAlgorithm = saml_conf["digest_algorithm"].(string)
	}
	if _, ok := saml_conf["default_relay_state"]; ok {
		res.DefaultRelayState = saml_conf["default_relay_state"].(string)
	}
	return res
}

func NewAppOidc(d *schema.ResourceData) *AppOidc {
	res := &AppOidc{}
	o, exists := d.GetOk("oidc")
	if !exists {
		return nil
	}
	oidc := o.([]interface{})
	if len(oidc) != 1 {
		return nil
	}
	oidc_conf := oidc[0].(map[string]interface{})
	res.SignInRedirectUrls = ResourceTypeSetToStringSlice(oidc_conf["sign_in_redirect_urls"].(*schema.Set))
	res.Scopes = ResourceTypeSetToStringSlice(oidc_conf["scopes"].(*schema.Set))
	if _, ok := oidc_conf["grant_types"]; ok {
		res.GrantTypes = ResourceTypeSetToStringSlice(oidc_conf["grant_types"].(*schema.Set))
	}
	ilu, ok := oidc_conf["initiate_login_url"]
	if ok && ilu != "" {
		ilu := ilu.(string)
		res.InitiateLoginUrl = &ilu
	} else {
		res.InitiateLoginUrl = nil
	}
	return res
}

func NewAppMappedAttr(d *schema.ResourceData) *[]AppMappedAttributes {
	ma, exists := d.GetOkExists("mapped_attributes")
	if !exists {
		return nil
	}
	rawMappedAttrs := ma.([]interface{})
	res := make([]AppMappedAttributes, len(rawMappedAttrs))
	for index, rawAttr := range rawMappedAttrs {
		rawAttrConf := rawAttr.(map[string]interface{})
		res[index].VariableName = rawAttrConf["variable_name"].(string)
		af, ok := rawAttrConf["attribute_format"]
		if ok && af != "" {
			af := af.(string)
			res[index].AttributeFormat = &af
		}
		tvn, ok := rawAttrConf["target_variable_name"]
		if ok && tvn != "" {
			tvn := tvn.(string)
			res[index].TargetVariableName = &tvn
		} else {
			res[index].TargetVariableName = nil
		}
		ft, ok := rawAttrConf["filter_type"]
		if ok && ft != "" {
			ft := ft.(string)
			res[index].FilterType = &ft
		} else {
			res[index].FilterType = nil
		}
		fv, ok := rawAttrConf["filter_value"]
		if ok && fv != "" {
			fv := fv.(string)
			res[index].FilterValue = &fv
		} else {
			res[index].FilterValue = nil
		}
	}
	return &res
}

func NewAppDomainFederation(protocol string, d *schema.ResourceData) (*AppDomainFederation, error) {
	res := &AppDomainFederation{}
	df, exists := d.GetOk("domain_federation")
	if !exists {
		return nil, nil
	}
	if protocol != "SAML" {
		return nil, fmt.Errorf("Domain federation with sso protocol %s is not supported", protocol)
	}
	domain_fed := df.([]interface{})
	if len(domain_fed) != 1 {
		return nil, nil
	}
	domain_federation_conf := domain_fed[0].(map[string]interface{})
	res.Domain = domain_federation_conf["domain"].(string)
	return res, nil
}

func parseApp(resp []byte) (*App, error) {
	app := &App{}
	err := json.Unmarshal(resp, app)
	if err != nil {
		return nil, fmt.Errorf("could not parse app response: %v", err)
	}
	return app, nil
}

func parseAppSaml(resp []byte) (*AppSaml, error) {
	app_saml := &AppSaml{}
	err := json.Unmarshal(resp, app_saml)
	if err != nil {
		return nil, fmt.Errorf("could not parse app saml response: %v", err)
	}
	return app_saml, nil
}

func parseAppOidc(resp []byte) (*AppOidc, error) {
	app_oidc := &AppOidc{}
	err := json.Unmarshal(resp, app_oidc)
	if err != nil {
		return nil, fmt.Errorf("could not parse app oidc response: %v", err)
	}
	return app_oidc, nil
}

func parseAppMappedAttributes(resp []byte) ([]AppMappedAttributes, error) {
	app_mapped_attrs := &[]AppMappedAttributes{}
	err := json.Unmarshal(resp, app_mapped_attrs)
	if err != nil {
		return nil, fmt.Errorf("could not parse app mapped attrs response: %v", err)
	}
	return *app_mapped_attrs, nil
}

func parseAppDomainFederation(resp []byte) (*AppDomainFederation, error) {
	app_domain_federation := &AppDomainFederation{}
	err := json.Unmarshal(resp, app_domain_federation)
	if err != nil {
		return nil, fmt.Errorf("could not parse app domain federation response: %v", err)
	}
	return app_domain_federation, nil
}

func UpdateAppProto(ctx context.Context, c *Client, app *App, saml []byte, oidc []byte,
	delete_on_failure bool) (*App, error) {
	if saml != nil {
		samlUrl := fmt.Sprintf("%s/%s/%s/saml", c.BaseURL, appEndpoint, app.ID)
		resp, err := c.Patch(ctx, samlUrl, saml)
		if err != nil {
			if delete_on_failure {
				DeleteApp(ctx, c, app.ID)
			}
			return nil, err
		}
		saml_resp, err := parseAppSaml(resp)
		if err != nil {
			if delete_on_failure {
				DeleteApp(ctx, c, app.ID)
			}
			return nil, err
		}
		app.Saml = saml_resp
	} else if oidc != nil {
		oidcUrl := fmt.Sprintf("%s/%s/%s/oidc", c.BaseURL, appEndpoint, app.ID)
		resp, err := c.Patch(ctx, oidcUrl, oidc)
		if err != nil {
			if delete_on_failure {
				DeleteApp(ctx, c, app.ID)
			}
			return nil, err
		}
		oidc_resp, err := parseAppOidc(resp)
		if err != nil {
			if delete_on_failure {
				DeleteApp(ctx, c, app.ID)
			}
			return nil, err
		}
		app.Oidc = oidc_resp
	}
	return app, nil
}

func UpdateAppDomainFederation(ctx context.Context, c *Client, app_id string,
	domainFed []byte) (*AppDomainFederation, error) {
	DomainFedUrl := fmt.Sprintf("%s/%s/%s/domain_federation", c.BaseURL, appEndpoint, app_id)
	resp, err := c.Patch(ctx, DomainFedUrl, domainFed)
	if err != nil {
		return nil, err
	}
	domain_federation_resp, err := parseAppDomainFederation(resp)
	if err != nil {
		return nil, err
	}
	return domain_federation_resp, nil
}

func UpdateAppMappedAttrs(ctx context.Context, c *Client, app_id string,
	mappedAttrs []byte) ([]AppMappedAttributes, error) {
	MappedAttrsUrl := fmt.Sprintf("%s/%s/%s/attribute_mapping", c.BaseURL, appEndpoint, app_id)
	resp, err := c.Put(ctx, MappedAttrsUrl, mappedAttrs)
	if err != nil {
		return nil, err
	}
	mapped_attrs_resp, err := parseAppMappedAttributes(resp)
	if err != nil {
		return nil, err
	}
	return mapped_attrs_resp, nil
}

func MarshalAppProtocol(protocol string, saml *AppSaml, oidc *AppOidc) ([]byte, []byte, error) {
	var saml_body []byte
	var oidc_body []byte
	var err error
	if protocol == "SAML" {
		saml_body, err = json.Marshal(saml)
		if err != nil {
			return nil, nil, fmt.Errorf("could not convert app saml to json: %v", err)
		}
	} else if protocol == "OIDC" {
		oidc_body, err = json.Marshal(oidc)
		if err != nil {
			return nil, nil, fmt.Errorf("could not convert app oidc to json: %v", err)
		}
	}
	return saml_body, oidc_body, nil
}

func CreateApp(ctx context.Context, c *Client, app *App, saml *AppSaml, oidc *AppOidc,
	mappedAttrs *[]AppMappedAttributes, domainFed *AppDomainFederation) (*App, error) {
	body, err := json.Marshal(app)
	if err != nil {
		return nil, fmt.Errorf("could not convert app to json: %v", err)
	}
	saml_body, oidc_body, err := MarshalAppProtocol(app.Protocol, saml, oidc)
	if err != nil {
		return nil, err
	}
	appUrl := fmt.Sprintf("%s/%s", c.BaseURL, appEndpoint)
	resp, err := c.Post(ctx, appUrl, body)
	if err != nil {
		return nil, err
	}
	app_resp, err := parseApp(resp)
	if err != nil {
		return nil, err
	}
	if mappedAttrs != nil {
		mapped_attrs_body, err := json.Marshal(mappedAttrs)
		if err != nil {
			DeleteApp(ctx, c, app_resp.ID)
			return nil, fmt.Errorf("could not convert app mapped attributes to json: %v", err)
		}
		if mapped_attrs_body != nil {
			mappedAttrs, err := UpdateAppMappedAttrs(ctx, c, app_resp.ID, mapped_attrs_body)
			if err != nil {
				DeleteApp(ctx, c, app_resp.ID)
				return nil, err
			}
			app_resp.MappedAttributes = mappedAttrs
		}
	}
	if domainFed != nil {
		domain_federation_body, err := json.Marshal(domainFed)
		if err != nil {
			DeleteApp(ctx, c, app_resp.ID)
			return nil, fmt.Errorf("could not convert app domain federation to json: %v", err)
		}
		if domain_federation_body != nil {
			domainFed, err := UpdateAppDomainFederation(ctx, c, app_resp.ID, domain_federation_body)
			if err != nil {
				DeleteApp(ctx, c, app_resp.ID)
				return nil, err
			}
			app_resp.DomainFederation = domainFed
		}
	}
	return UpdateAppProto(ctx, c, app_resp, saml_body, oidc_body, true)
}

func UpdateApp(ctx context.Context, c *Client, appID string, app *App, saml *AppSaml, oidc *AppOidc,
	mappedAttrs *[]AppMappedAttributes, domainFed *AppDomainFederation) (*App, error) {
	var empty_proto string
	proto := app.Protocol
	app.Protocol = empty_proto

	body, err := json.Marshal(app)
	if err != nil {
		return nil, fmt.Errorf("could not convert app rule to json: %v", err)
	}
	saml_body, oidc_body, err := MarshalAppProtocol(proto, saml, oidc)
	if err != nil {
		return nil, err
	}
	appUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, appEndpoint, appID)
	resp, err := c.Patch(ctx, appUrl, body)
	if err != nil {
		return nil, err
	}
	app_resp, err := parseApp(resp)
	if err != nil {
		return nil, err
	}
	if mappedAttrs != nil {
		mapped_attrs_body, err := json.Marshal(mappedAttrs)
		if err != nil {
			return nil, fmt.Errorf("could not convert app mapped attributes to json: %v", err)
		}
		if mapped_attrs_body != nil {
			mappedAttrs, err := UpdateAppMappedAttrs(ctx, c, app_resp.ID, mapped_attrs_body)
			if err != nil {
				return nil, err
			}
			app_resp.MappedAttributes = mappedAttrs
		}
	}
	if domainFed != nil {
		domain_federation_body, err := json.Marshal(domainFed)
		if err != nil {
			return nil, fmt.Errorf("could not convert app domain federation to json: %v", err)
		}
		if domain_federation_body != nil {
			domainFed, err := UpdateAppDomainFederation(ctx, c, app_resp.ID, domain_federation_body)
			if err != nil {
				return nil, err
			}
			app_resp.DomainFederation = domainFed
		}
	}
	app.Protocol = proto
	return UpdateAppProto(ctx, c, app_resp, saml_body, oidc_body, false)
}

func GetApp(ctx context.Context, c *Client, appID string, protocol string) (*App, error) {
	appUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, appEndpoint, appID)
	resp, err := c.Get(ctx, appUrl, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	app_resp, err := parseApp(resp)
	if err != nil {
		return nil, err
	}
	if protocol == "SAML" {
		samlUrl := fmt.Sprintf("%s/%s/%s/saml", c.BaseURL, appEndpoint, app_resp.ID)
		resp, err = c.Get(ctx, samlUrl, nil)
		if err != nil {
			return nil, err
		}
		saml_resp, err := parseAppSaml(resp)
		if err != nil {
			return nil, err
		}
		app_resp.Saml = saml_resp
		DomainFedUrl := fmt.Sprintf("%s/%s/%s/domain_federation", c.BaseURL, appEndpoint, app_resp.ID)
		resp, err = c.Get(ctx, DomainFedUrl, nil)
		if err == nil {
			domain_federation_resp, err := parseAppDomainFederation(resp)
			if err != nil {
				return nil, err
			}
			app_resp.DomainFederation = domain_federation_resp
		}
	} else if protocol == "OIDC" {
		oidcUrl := fmt.Sprintf("%s/%s/%s/oidc", c.BaseURL, appEndpoint, app_resp.ID)
		resp, err = c.Get(ctx, oidcUrl, nil)
		if err != nil {
			return nil, err
		}
		oidc_resp, err := parseAppOidc(resp)
		if err != nil {
			return nil, err
		}
		app_resp.Oidc = oidc_resp
	}
	MappedAttrsUrl := fmt.Sprintf("%s/%s/%s/attribute_mapping", c.BaseURL, appEndpoint, app_resp.ID)
	resp, err = c.Get(ctx, MappedAttrsUrl, nil)
	if err != nil {
		return nil, err
	}
	mapped_attrs_resp, err := parseAppMappedAttributes(resp)
	if err != nil {
		return nil, err
	}
	app_resp.MappedAttributes = mapped_attrs_resp
	return app_resp, nil
}

func DeleteApp(ctx context.Context, c *Client, appID string) (*App, error) {
	appUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, appEndpoint, appID)
	resp, err := c.Delete(ctx, appUrl, nil)
	if err != nil {
		return nil, err
	}
	return parseApp(resp)
}
