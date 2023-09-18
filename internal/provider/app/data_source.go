package app

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: appRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "app"),
			},
			"protocol": {
				Description:      protoDesc,
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateStringENUM("SAML", "OIDC"),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"visible": {
				Description: visibleDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"assigned_members": {
				Description: assignedMembersDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"ip_whitelist": {
				Description: assignedMembersDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"direct_sso_login": {
				Description: directSsoLoginDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"saml": {
				Description: samlDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audience_uri": {
							Description: samlAudienceUriDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"recipient": {
							Description: samlRecipientDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"destination": {
							Description: samlDestinationDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"sso_acs_url": {
							Description: samlSsoAcsUrlDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"subject_name_id_attribute": {
							Description: samlSubNameIDAttrDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"subject_name_id_format": {
							Description: samlSubNameIDFormatDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"signature_algorithm": {
							Description: samlSigAlgDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"digest_algorithm": {
							Description: samlDigestAlgDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"default_relay_state": {
							Description: samlDefRelayStateDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"x509_cert": {
							Description: samlx509CertDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"idp_issuer": {
							Description: samlIdpIssuerDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"idp_sso_url": {
							Description: samlSsoUrleDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"authn_context_class": {
							Description: samlAuthnCtxClassDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"oidc": {
				Description: oidcDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sign_in_redirect_urls": {
							Description: oidcSigninRedUrlsDesc,
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
						},
						"grant_types": {
							Description: oidcGrantTypesDesc,
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
						},
						"scopes": {
							Description: oidcScopesDesc,
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
						},
						"initiate_login_url": {
							Description: oidcInitLoginUrlDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"access_token_lifetime": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"id_token_lifetime": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"mapped_attributes": {
				Description: MappedAttributesDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attribute_format": {
							Description: attrFmtDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"target_variable_name": {
							Description: attrTargetNameDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"variable_name": {
							Description: attrVarNameDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"filter_type": {
							Description: attrFilterTypeDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"filter_value": {
							Description: attrFilterValueDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
