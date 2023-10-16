package app

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: appCreate,
		ReadContext:   appRead,
		UpdateContext: appUpdate,
		DeleteContext: appDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"visible": {
				Description: visibleDesc,
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"assigned_members": {
				Description: assignedMembersDesc,
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "usr", "grp")},
			},
			"ip_whitelist": {
				Description: ipWhiteListDesc,
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type: schema.TypeString},
			},
			"direct_sso_login": {
				Description:      directSsoLoginDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateID(false, "idp"),
			},
			"protocol": {
				Description:      protoDesc,
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateStringENUM("SAML", "OIDC"),
			},
			"saml": {
				Description: samlDesc,
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audience_uri": {
							Description: samlAudienceUriDesc,
							Type:        schema.TypeString,
							Required:    true,
						},
						"recipient": {
							Description: samlRecipientDesc,
							Type:        schema.TypeString,
							Required:    true,
						},
						"destination": {
							Description: samlDestinationDesc,
							Type:        schema.TypeString,
							Required:    true,
						},
						"sso_acs_url": {
							Description: samlSsoAcsUrlDesc,
							Type:        schema.TypeString,
							Required:    true,
						},
						"subject_name_id_attribute": {
							Description:      samlSubNameIDAttrDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateStringENUM("email", "immutableid"),
						},
						"subject_name_id_format": {
							Description:      samlSubNameIDFormatDesc,
							Type:             schema.TypeString,
							Default:          "unspecified",
							Optional:         true,
							ValidateDiagFunc: common.ValidateStringENUM("unspecified", "emailAddress", "persistent"),
						},
						"signature_algorithm": {
							Description:      samlSigAlgDesc,
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "RSA-SHA256",
							ValidateDiagFunc: common.ValidateStringENUM("RSA-SHA256"),
						},
						"digest_algorithm": {
							Description:      samlDigestAlgDesc,
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "SHA256",
							ValidateDiagFunc: common.ValidateStringENUM("SHA256"),
						},
						"default_relay_state": {
							Description: samlDefRelayStateDesc,
							Type:        schema.TypeString,
							Optional:    true,
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
				ConflictsWith: []string{"oidc"},
			},
			"domain_federation": {
				Description: domainFedDesc,
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"domain": {
							Description:      domainFedDomainDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidatePattern(common.DomainPattern),
						},
					},
				},
				ConflictsWith: []string{"oidc"},
			},
			"oidc": {
				Description: oidcDesc,
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sign_in_redirect_urls": {
							Description: oidcSigninRedUrlsDesc,
							Type:        schema.TypeSet,
							Required:    true,
							MinItems:    1,
							MaxItems:    5,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidateURL()},
						},
						"grant_types": {
							Description: oidcGrantTypesDesc,
							Type:        schema.TypeSet,
							Optional:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidateStringENUM("authorization_code")},
						},
						"scopes": {
							Description: oidcScopesDesc,
							Type:        schema.TypeSet,
							Required:    true,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidateStringENUM("openid", "email", "profile", "phone", "groups")},
						},
						"initiate_login_url": {
							Description:      oidcInitLoginUrlDesc,
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: common.ValidateURL(),
						},
						"access_token_lifetime": {
							Description: oidcAccessTokenLifetimeDesc,
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"id_token_lifetime": {
							Description: oidcIdTokenLifetimeDesc,
							Type:        schema.TypeInt,
							Computed:    true,
						},
					},
				},
				ConflictsWith: []string{"saml", "domain_federation"},
			},
			"mapped_attributes": {
				Description: MappedAttributesDesc,
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    15,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"attribute_format": {
							Description:      attrFmtDesc,
							Type:             schema.TypeString,
							Optional:         true,
							Default:          "unspecified",
							ValidateDiagFunc: common.ValidateStringENUM("unspecified", "basic", "uri"),
						},
						"target_variable_name": {
							Description: attrTargetNameDesc,
							Type:        schema.TypeString,
							Optional:    true,
						},
						"variable_name": {
							Description: attrVarNameDesc,
							Type:        schema.TypeString,
							Required:    true,
							ValidateDiagFunc: common.ValidateStringENUM("given_name", "family_name", "email",
								"phone", "groups", "tags"),
						},
						"filter_type": {
							Description:      attrFilterTypeDesc,
							Type:             schema.TypeString,
							Optional:         true,
							ValidateDiagFunc: common.ValidateStringENUM("equals", "all", "starts_with"),
						},
						"filter_value": {
							Description: attrFilterValueDesc,
							Type:        schema.TypeString,
							Optional:    true,
						},
					},
				},
			},
		},
	}
}
