package dlp_rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

const maxInt = int(^uint(0) >> 1)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: dlpRuleCreate,
		ReadContext:   dlpRuleRead,
		UpdateContext: dlpRuleUpdate,
		DeleteContext: dlpRuleDelete,
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
				ValidateDiagFunc: common.ValidateStringENUM("BLOCK", "LOG"),
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
					ValidateDiagFunc: common.ValidateID(false, "usr", "grp"),
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
					ValidateDiagFunc: common.ValidateID(false, "usr", "grp"),
				},
				Optional:      true,
				ConflictsWith: []string{"apply_to_org"},
			},
			"alert_level": {
				Description:      alertLevelDesc,
				Type:             schema.TypeString,
				ValidateDiagFunc: common.ValidateStringENUM("CRITICAL", "HIGH", "MEDIUM", "LOW"),
				Optional:         true,
			},
			"all_resources": {
				Description:   allResourcesDesc,
				Type:          schema.TypeBool,
				ConflictsWith: []string{"resource_countries", "cloud_apps", "content_types", "threat_types"},
				Optional:      true,
			},
			"all_supported_file_types": {
				Description:   allSupportedFileTypeDesc,
				Type:          schema.TypeBool,
				ConflictsWith: []string{"file_types"},
				Optional:      true,
			},
			"cloud_apps": {
				Description: cloudAppsDesc,
				Type:        schema.TypeList,
				MaxItems:    50,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "ca"),
				},
				Optional:      true,
				ConflictsWith: []string{"all_resources"},
			},
			"content_types": {
				Description: contentTypesDesc,
				Type:        schema.TypeList,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM(common.ContentTypes...),
				},
				Optional:      true,
				ConflictsWith: []string{"all_resources"},
			},
			"detectors": {
				Description: detectorsDesc,
				Type:        schema.TypeList,
				MaxItems:    100,
				MinItems:    1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"file_parts": {
				Description: filePartsDesc,
				Type:        schema.TypeList,
				MaxItems:    3,
				MinItems:    1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("CONTENT", "FILE_NAME", "METADATA"),
				},
				Required: true,
			},
			"file_types": {
				Description: fileTypesDesc,
				Type:        schema.TypeList,
				MaxItems:    3,
				MinItems:    1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("arj", "csv", "doc", "docm", "docx", "dot",
						"dotm", "dotx", "java", "odp", "ods", "odt", "pdf", "perl", "pot", "potm", "potx", "ppa",
						"ppam", "pps", "ppsm", "ppsx", "ppt", "pptm", "pptx", "py", "rtf", "sh", "sql", "txt", "xla",
						"xlam", "xlm", "xls", "xlsb", "xlsm", "xlsx", "xlt", "xltm", "xltx", "zip"),
				},
				ConflictsWith: []string{"all_supported_file_types"},
				Optional:      true,
			},
			"filter_expression": {
				Description: expressionDesc,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"priority": {
				Description:      priorityDesc,
				Type:             schema.TypeInt,
				ValidateDiagFunc: common.ValidateIntRange(1, 5000),
				Required:         true,
			},
			"resource_countries": {
				Description: resourceCountriesDesc,
				Type:        schema.TypeList,
				MaxItems:    10,
				MinItems:    1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM(common.Countries...),
				},
				ConflictsWith: []string{"all_resources"},
				Optional:      true,
			},
			"threat_types": {
				Description: threatTypesDesc,
				Type:        schema.TypeList,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("Abused TLD", "Bitcoin Related", "Blackhole",
						"Botnets", "Brute Forcer", "Chat Server", "CnC", "Compromised", "DDoS Target", "Drop",
						"DynDNS", "EXE Source", "Fake AV", "IP Check", "Keyloggers and Monitoring", "Malware Sites",
						"Mobile CnC", "Mobile Spyware CnC", "Online Gaming", "P2P CnC", "Peer to Peer", "Parking",
						"Phishing and Other Frauds", "Private IP Addresses", "Proxy Avoidance and Anonymizers",
						"Remote Access Service", "Scanner", "Self Signed SSL", "SPAM URLs", "Spyware and Adware",
						"Tor", "Undesirable", "Utility", "VPN"),
				},
				ConflictsWith: []string{"all_resources"},
				Optional:      true,
			},
			"user_actions": {
				Description: userActionsDesc,
				Type:        schema.TypeList,
				MinItems:    1,
				MaxItems:    8,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("DOWNLOAD", "UPLOAD", "CREATE", "EDIT", "SHARE",
						"POST", "DELETE", "LIKE"),
				},
				Required: true,
			},
		},
	}
}
