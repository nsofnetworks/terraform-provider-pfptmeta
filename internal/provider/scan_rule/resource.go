package scan_rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: scanRuleCreate,
		ReadContext:   scanRuleRead,
		UpdateContext: scanRuleUpdate,
		DeleteContext: scanRuleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"priority": {
				Description:      priorityDesc,
				Type:             schema.TypeInt,
				ValidateDiagFunc: common.ValidateIntRange(1, 5000),
				Required:         true,
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
				Description: exemptSourcesDesc,
				Type:        schema.TypeList,
				MaxItems:    200,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "usr", "grp"),
				},
				Optional:      true,
				ConflictsWith: []string{"apply_to_org"},
			},
			"filter_expression": {
				Description: expressionDesc,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"apply_to_org": {
				Description:   applyToOrgDesc,
				Type:          schema.TypeBool,
				Optional:      true,
				ConflictsWith: []string{"sources", "exempt_sources"},
			},
			"content_categories": {
				Description: contentCategoriesDesc,
				Type:        schema.TypeList,
				MaxItems:    20,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "cc"),
				},
				Optional: true,
			},
			"threat_categories": {
				Description: threatCategoriesDesc,
				Type:        schema.TypeList,
				MaxItems:    5,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "tc"),
				},
				Optional: true,
			},
			"cloud_apps": {
				Description: cloudAppsDesc,
				Type:        schema.TypeList,
				MaxItems:    50,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "ca"),
				},
				Optional: true,
			},
			"cloud_app_risk_groups": {
				Description: cloudAppRiskGroupsDesc,
				Type:        schema.TypeList,
				MaxItems:    10,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("imminent_targets", "latent_targets", "major_targets",
						"soft_targets", "very_attacked_apps", "very_privileged_apps",
						"very_vulnerable_apps"),
				},
				Optional: true,
			},
			"catalog_app_categories": {
				Description: catalogAppCategoriesDesc,
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
			"networks": {
				Description: networksDesc,
				Type:        schema.TypeList,
				MaxItems:    20,
				MinItems:    1,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "ipn"),
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
			"user_agents": {
				Description: userAgentsDesc,
				Type:        schema.TypeList,
				MinItems:    1,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("Chrome", "Safari", "Edge",
						"Firefox", "Opera", "IE", "Electron",
						"Outlook", "Excel", "PowerPoint"),
				},
				Optional: true,
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
			"all_supported_file_types": {
				Description:   allSupportedFileTypeDesc,
				Type:          schema.TypeBool,
				ConflictsWith: []string{"file_types"},
				Optional:      true,
			},
			"file_types": {
				Description: fileTypesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("7z", "ace", "arj", "bat", "cab", "chm", "contact", "cpl",
						"csv", "dll", "dmg", "doc", "docb", "docm", "docx", "dot",
						"dotm", "dotx", "exe", "form", "gz", "gzip", "hta", "htm",
						"html", "img", "inetloc", "iqy", "iso", "jar", "java", "jnlp",
						"js", "lnk", "mam", "mht", "mhtml", "msi", "odp", "ods", "odt",
						"one", "onepkg", "pdf", "perl", "php", "pkg", "plist", "pot",
						"potm", "potx", "ppa", "ppam", "pps", "ppsm", "ppsx", "ppt",
						"pptm", "pptx", "ps", "ps1", "ps1xml", "ps2", "ps2xml", "psc1",
						"psc2", "pub", "py", "rar", "reg", "rtf", "sh", "shtm", "shtml",
						"slk", "sql", "swf", "txt", "vb", "vbe", "vbs", "vbscript", "vcard",
						"vcf", "vcs", "vhd", "vhdx", "wll", "wmv", "xht", "xla", "xlam", "xll",
						"xlm", "xls", "xlsb", "xlsm", "xlsx", "xlt", "xltm", "xltx", "xps", "xxe", "zip"),
				},
				ConflictsWith: []string{"all_supported_file_types"},
				Optional:      true,
			},
			"max_file_size_mb": {
				Description:      maxFileSizeMBDesc,
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: common.ValidateIntRange(1, 4096),
			},
			"password_protected_files": {
				Description: passwordProtectedFilesDesc,
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"dlp": {
				Description: dlpDesc,
				Type:        schema.TypeBool,
				Optional:    true,
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
			"malware": {
				Description: malwareDesc,
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"sandbox": {
				Description: sandboxDesc,
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"antivirus": {
				Description: antivirusDesc,
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"action": {
				Description:      actionDesc,
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateStringENUM("BLOCK", "LOG"),
			},
		},
	}
}
