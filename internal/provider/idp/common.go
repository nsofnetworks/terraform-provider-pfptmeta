package idp

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

var excludedKeys = []string{"id", "saml_config", "oidc_config", "scim_config"}

const (
	description = "IdP configuration to allow logging into the Proofpoint NaaS platform via SSO using a 3rd-party IdP. " +
		"It is also used to configure users and groups provisioning from the 3rd-party IdP using SCIM protocol."
	hiddenDesc            = "Defines whether to display the IdP on Proofpoint login page as an option for SSO login"
	iconDesc              = "Icon to display on Proofpoint login page"
	mappedAttributesDesc  = "User attributes to map from IdP to Proofpoint platform. It can be provided using SSO-JIT or SCIM"
	samlDesc              = "SSO configuration using SAML protocol"
	samlSsoUrlDesc        = "3rd-party IdP SSO login URL"
	samlAuthnCtxClassDesc = "Authentication method to request the 3rd-part IdP to comply with during SAML request. Options: " +
		"unspecified, Password, PasswordProtectedTransport, X509, Smartcard, Kerberos"
	samlCertDesc         = "3rd-Party IdP certificate to be used by SAML for signature validation"
	jitEnabledDesc       = "Defines whether to allow just-in-time user provisioning during SSO login"
	issuerDesc           = "3rd-party IdP issuer"
	oidcDesc             = "SSO configuration using OIDC protocol"
	oidcClientIdDesc     = "Client ID as published by the 3rd-party IdP"
	oidcClientSecretDesc = "Client secret as published by 3rd-party IdP"
	scimDesc             = "Provisioning configuration using SCIM protocol"
	scimApiKeyIdDesc     = "API key ID to be used by the 3rd-party IdP to make API calls to Proofpoint platform. " +
		"It is also used for user and group provisioning"
	scimAssumeOwnershipDesc = "Defines whether to take ownership over resources that are not provisioned"
)

func idpRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	i, err := client.GetIdp(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing Idp %s because it's gone", id)
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	return idpToResource(d, i)
}

func idpCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewIdp(d)
	i, err := client.CreateIdp(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return idpToResource(d, i)
}

func idpUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewIdp(d)
	i, err := client.UpdateIdp(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return idpToResource(d, i)
}

func idpDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteIdp(ctx, c, id)
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

func idpToResource(d *schema.ResourceData, i *client.Idp) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId(i.ID)
	err := client.MapResponseToResource(i, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	if i.SamlConfig != nil {
		samlToResource := []map[string]interface{}{
			{"certificate": i.SamlConfig.Certificate, "issuer": i.SamlConfig.Issuer,
				"sso_url": i.SamlConfig.SsoUrl, "authn_context_class": i.SamlConfig.AuthnContextClass,
				"jit_enabled": i.SamlConfig.JitEnabled},
		}
		err = d.Set("saml_config", samlToResource)
		if err != nil {
			return diag.FromErr(err)
		}
	} else if i.OidcConfig != nil {
		oidcToResource := []map[string]interface{}{
			{"issuer": i.OidcConfig.Issuer, "client_id": i.OidcConfig.ClientId,
				"jit_enabled": i.OidcConfig.JitEnabled},
		}
		err = d.Set("oidc_config", oidcToResource)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if i.ScimConfig != nil {
		ScimToResource := []map[string]interface{}{
			{"api_key_id": i.ScimConfig.ApiKeyId, "assume_ownership": i.ScimConfig.AssumeOwnership},
		}
		err = d.Set("scim_config", ScimToResource)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}
