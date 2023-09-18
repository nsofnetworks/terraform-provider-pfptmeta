package app

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

var excludedKeys = []string{"id", "saml", "oidc", "mapped_attributes"}

const (
	description         = "Application for configuring SSO by SPs based on SAML or OIDC protocols"
	visibleDesc         = "Application visibility, defining whether to display application to user or not"
	assignedMembersDesc = "Users and groups which the application is applied to"
	ipWhiteListDesc     = "List of IPs allowed to be authenticated by the application"
	directSsoLoginDesc  = "IDP to use when logging into Proofpoint NaaS directly, while performing SSO login to SP"
	protoDesc           = "The protocol that the SSO uses for this application. After defining the protocol on SSO creation, " +
		"it cannot be changed. Options: SAML/OIDC"
	samlDesc                    = "SAML-based app properties"
	samlAudienceUriDesc         = "SP Entity ID"
	samlRecipientDesc           = "The location for the app to present the SAML assertion"
	samlDestinationDesc         = "Destination within the SAML assertion"
	samlSsoAcsUrlDesc           = "Single Sign-On URL"
	samlSubNameIDAttrDesc       = "SAML name ID type to identify the subject of the SSO request"
	samlSubNameIDFormatDesc     = "ID format for SAML name"
	samlSigAlgDesc              = "SAML assertion signature algorithm"
	samlDigestAlgDesc           = "SAML assertion digest algorithm"
	samlx509CertDesc            = "SAML certificate to be configured at SP side"
	samlIdpIssuerDesc           = "SAML issuer to be configured at the SP side"
	samlSsoUrleDesc             = "SAML url to be configured at SP side"
	samlAuthnCtxClassDesc       = "SAML authentication context class to be configured at the SP side"
	samlDefRelayStateDesc       = "SAML default relay state URL to use after successful assertion"
	wsFedDesc                   = "SSO configuration for Office365 SP"
	oidcDesc                    = "OIDC-based app properties"
	oidcSigninRedUrlsDesc       = "Redirect URLs which are allowed after successful authorization"
	oidcGrantTypesDesc          = "OIDC-supported access/ID token grant types"
	oidcScopesDesc              = "Scopes of allowed access token for OIDC"
	oidcInitLoginUrlDesc        = "Login URL for IdP-initiated flow"
	oidcAccessTokenLifetimeDesc = "Validity period (in minutes) for access token from the moment it is generated"
	oidcIdTokenLifetimeDesc     = "Validity period (in minutes) for ID token from the moment it is generated"
	MappedAttributesDesc        = "User attributes to map and return to SP upon successful SAML assertion/OIDC authorization"
	attrFmtDesc                 = "Format in which the name attribute is provided to the app. Options: unspecified, basic, uri"
	attrTargetNameDesc          = "Name for the attribute as it is presented in SP"
	attrVarNameDesc             = "Variable name where the value for the attribute is acquired from. " +
		"Options: given_name, family_name, email, phone, groups, tags"
	attrFilterTypeDesc  = "Filter type to use when returning the attributes to SP. Options: all, starts_with, equals"
	attrFilterValueDesc = "Filter type value ('starts_with' and 'equals') to compare the actual value to"
)

func appRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Get("id").(string)
	proto := d.Get("protocol").(string)
	c := meta.(*client.Client)
	a, err := client.GetApp(ctx, c, id, proto)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing application %s because it's gone", id)
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	return appToResource(d, a)
}

func appCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	app_body := client.NewApp(d)
	var saml_body *client.AppSaml
	var oidc_body *client.AppOidc
	if app_body.Protocol == "SAML" {
		saml_body = client.NewAppSaml(d)
	} else if app_body.Protocol == "OIDC" {
		oidc_body = client.NewAppOidc(d)
	}
	mapped_attrs_body := client.NewAppMappedAttr(d)
	a, err := client.CreateApp(ctx, c, app_body, saml_body, oidc_body, mapped_attrs_body)
	if err != nil {
		return diag.FromErr(err)
	}
	return appToResource(d, a)
}

func appUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	app_body := client.NewApp(d)
	var saml_body *client.AppSaml
	var oidc_body *client.AppOidc
	if app_body.Protocol == "SAML" {
		saml_body = client.NewAppSaml(d)
	} else if app_body.Protocol == "OIDC" {
		oidc_body = client.NewAppOidc(d)
	}
	mapped_attrs_body := client.NewAppMappedAttr(d)
	a, err := client.UpdateApp(ctx, c, id, app_body, saml_body, oidc_body, mapped_attrs_body)
	if err != nil {
		return diag.FromErr(err)
	}
	return appToResource(d, a)
}

func appDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteApp(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	return
}

func appToResource(d *schema.ResourceData, a *client.App) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId(a.ID)
	err := client.MapResponseToResource(a, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	if a.Saml != nil {
		samlToResource := []map[string]interface{}{
			{"audience_uri": a.Saml.AudienceUri, "recipient": a.Saml.Recipient, "destination": a.Saml.Destination,
				"sso_acs_url": a.Saml.SsoAcsUrl, "subject_name_id_attribute": a.Saml.SubjectNameIdAttribute,
				"subject_name_id_format": a.Saml.SubjectNameIdFormat, "signature_algorithm": a.Saml.SignatureAlgorithm,
				"digest_algorithm": a.Saml.DigestAlgorithm, "default_relay_state": a.Saml.DefaultRelayState,
				"x509_cert": a.Saml.X509Cert, "idp_issuer": a.Saml.IdpIssuer, "idp_sso_url": a.Saml.IdpSsoUrl,
				"authn_context_class": a.Saml.AuthnContextClass},
		}
		err = d.Set("saml", samlToResource)
		if err != nil {
			return diag.FromErr(err)
		}
	} else if a.Oidc != nil {
		oidcToResource := []map[string]interface{}{
			{"sign_in_redirect_urls": a.Oidc.SignInRedirectUrls, "grant_types": a.Oidc.GrantTypes, "scopes": a.Oidc.Scopes,
				"initiate_login_url": a.Oidc.InitiateLoginUrl, "access_token_lifetime": a.Oidc.AccessTokenLifetime,
				"id_token_lifetime": a.Oidc.IdTokenLifetime},
		}
		err = d.Set("oidc", oidcToResource)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	var mappedAttrsToResource []map[string]interface{}
	for _, mappedAttr := range a.MappedAttributes {
		mappedAttr := map[string]interface{}{
			"attribute_format": mappedAttr.AttributeFormat, "variable_name": mappedAttr.VariableName,
			"target_variable_name": mappedAttr.TargetVariableName, "filter_type": mappedAttr.FilterType,
			"filter_value": mappedAttr.FilterValue}
		mappedAttrsToResource = append(mappedAttrsToResource, mappedAttr)
	}
	err = d.Set("mapped_attributes", mappedAttrsToResource)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
