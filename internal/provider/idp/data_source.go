package idp

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: idpRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "idp"),
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
			"hidden": {
				Description: hiddenDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"icon": {
				Description: iconDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"mapped_attributes": {
				Description: mappedAttributesDesc,
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"saml_config": {
				Description: samlDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"certificate": {
							Description: samlCertDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"issuer": {
							Description: issuerDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"sso_url": {
							Description: samlSsoUrlDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"authn_context_class": {
							Description: samlAuthnCtxClassDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"jit_enabled": {
							Description: jitEnabledDesc,
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
			"oidc_config": {
				Description: oidcDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"issuer": {
							Description: issuerDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"client_id": {
							Description: oidcClientIdDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"jit_enabled": {
							Description: jitEnabledDesc,
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
			"scim_config": {
				Description: scimDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"api_key_id": {
							Description: scimApiKeyIdDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"assume_ownership": {
							Description: scimAssumeOwnershipDesc,
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}
