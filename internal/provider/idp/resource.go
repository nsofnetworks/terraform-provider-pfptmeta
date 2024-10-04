package idp

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: idpCreate,
		ReadContext:   idpRead,
		UpdateContext: idpUpdate,
		DeleteContext: idpDelete,
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
			"hidden": {
				Description: hiddenDesc,
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"icon": {
				Description: iconDesc,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"mapped_attributes": {
				Description: mappedAttributesDesc,
				Type:        schema.TypeSet,
				Optional:    true,
				MaxItems:    10,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidatePattern(common.TagPattern)},
			},
			"saml_config": {
				Description: samlDesc,
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate": {
							Description: samlCertDesc,
							Type:        schema.TypeString,
							Required:    true,
						},
						"issuer": {
							Description: issuerDesc,
							Type:        schema.TypeString,
							Required:    true,
						},
						"sso_url": {
							Description:      samlSsoUrlDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateURL(),
						},
						"authn_context_class": {
							Description: samlAuthnCtxClassDesc,
							Type:        schema.TypeString,
							Default:     "PasswordProtectedTransport",
							Optional:    true,
							ValidateDiagFunc: common.ValidateStringENUM(
								"unspecified", "Password", "PasswordProtectedTransport", "X509",
								"Smartcard", "Kerberos"),
						},
						"jit_enabled": {
							Description: jitEnabledDesc,
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
					},
				},
				ConflictsWith: []string{"oidc_config"},
			},
			"oidc_config": {
				Description: oidcDesc,
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"issuer": {
							Description:      issuerDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateURL(),
						},
						"client_id": {
							Description: oidcClientIdDesc,
							Type:        schema.TypeString,
							Required:    true,
						},
						"client_secret": {
							Description: oidcClientSecretDesc,
							Type:        schema.TypeString,
							Required:    true,
							Sensitive:   true,
						},
						"jit_enabled": {
							Description: jitEnabledDesc,
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
						},
					},
				},
				ConflictsWith: []string{"saml_config"},
			},
			"scim_config": {
				Description: scimDesc,
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key_id": {
							Description:      scimApiKeyIdDesc,
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: common.ValidateID(false, "key"),
						},
						"assume_ownership": {
							Description: scimAssumeOwnershipDesc,
							Type:        schema.TypeBool,
							Required:    true,
						},
					},
				},
			},
		},
	}
}
