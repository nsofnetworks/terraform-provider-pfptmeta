package url_filtering_rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

const maxInt = int(^uint(0) >> 1)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: urlFilteringRuleCreate,
		ReadContext:   urlFilteringRuleRead,
		UpdateContext: urlFilteringRuleUpdate,
		DeleteContext: urlFilteringRuleDelete,
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
			"action": {
				Description:      actionDesc,
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateStringENUM("ISOLATION", "BLOCK", "LOG", "RESTRICT", "WARN"),
			},
			"apply_to_org": {
				Description:   applyToOrgDesc,
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"sources", "exempt_sources"},
			},
			"sources": {
				Description: sourcesDesc,
				Type:        schema.TypeList,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "usr", "grp", "tun"),
				},
				Optional:      true,
				ConflictsWith: []string{"apply_to_org"},
			},
			"exempt_sources": {
				Description: exemptSources,
				Type:        schema.TypeList,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "usr", "grp", "tun"),
				},
				Optional:      true,
				ConflictsWith: []string{"apply_to_org"},
			},
			"advanced_threat_protection": {
				Description: advancedThreatProtectionDesc,
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"catalog_app_categories": {
				Description: catalogAppCategories,
				Type:        schema.TypeList,
				MaxItems:    20,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("Instant Messaging", "eCommerce",
						"Content Management", "Software Development", "Project Management", "Marketing", "CRM",
						"Telecommunications", "Social and Communication", "Productivity", "Collaboration",
						"Business and Finance", "Utilities", "IT Service Management", "Social Networking",
						"Office Document and Productivity", "Cloud File Sharing", "Web Meetings",
						"Identity and Access Management", "IT Services and Hosting", "Webmail", "Website Builder",
						"Human Capital Management", "Sales and CRM", "E-commerce and Accounting", "Streaming Media",
						"Cloud Storage", "Operations Management", "Online Meeting", "Supply Chain",
						"Security and Compliance", "Entertainment and Lifestyle", "System and Network",
						"Retail and Consumer Services", "Health and Benefits", "Data and Analytics",
						"Education and References", "Personal instant messaging", "Legal", "Other",
						"Hosting Services", "News and Media", "Sales", "Enterprise Resource Planning", "Advertising",
						"Travel and Transportation", "Property Management", "Government Services", "Games", "Code Hosting"),
				},
				Optional: true,
			},
			"catalog_app_risk": {
				Description:      catalogAppRiskDesc,
				Type:             schema.TypeInt,
				ValidateDiagFunc: common.ValidateIntRange(1, 8),
				Optional:         true,
			},
			"cloud_apps": {
				Description:   cloudAppsDesc,
				Type:          schema.TypeList,
				MaxItems:      50,
				ConflictsWith: []string{"tenant_restriction", "tenant_restriction", "threat_categories"},
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "ca"),
				},
				Optional: true,
			},
			"countries": {
				Description: countriesDesc,
				Type:        schema.TypeList,
				MaxItems:    10,
				MinItems:    1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM(common.Countries...),
				},
				Optional: true,
			},
			"expires_at": {
				Description:      expiresAtDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateIsoTimeFormat(),
			},
			"filter_expression": {
				Description: expressionDesc,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"forbidden_content_categories": {
				Description:   contentCategoriesDesc,
				Type:          schema.TypeList,
				MaxItems:      20,
				ConflictsWith: []string{"tenant_restriction", "cloud_apps"},
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "cc"),
				},
				Optional: true,
			},
			"networks": {
				Description: networkDesc,
				Type:        schema.TypeList,
				MaxItems:    20,
				MinItems:    1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "ipn"),
				},
				Optional: true,
			},
			"priority": {
				Description:      priorityDesc,
				Type:             schema.TypeInt,
				ValidateDiagFunc: common.ValidateIntRange(1, 5000),
				Optional:         true,
			},
			"schedule": {
				Description: scheduleDesc,
				Type:        schema.TypeList,
				MaxItems:    10,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "tmf"),
				},
				Optional: true,
			},
			"tenant_restriction": {
				Description:      tenantRestrictionDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateID(false, "tr"),
				ConflictsWith:    []string{"forbidden_content_categories", "cloud_apps", "threat_categories"},
			},
			"threat_categories": {
				Description:   threatCategoriesDesc,
				Type:          schema.TypeList,
				ConflictsWith: []string{"tenant_restriction", "cloud_apps"},
				MaxItems:      5,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "tc"),
				},
				Optional: true,
			},
			"warn_ttl": {
				Description:      warnTtlDesc,
				Type:             schema.TypeInt,
				ValidateDiagFunc: common.ValidateIntRange(1, 43800),
				Optional:         true,
			},
		},
	}
}
