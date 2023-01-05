package file_scanning_rule

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

const maxInt = int(^uint(0) >> 1)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: fileScanningRuleCreate,
		ReadContext:   fileScanningRuleRead,
		UpdateContext: fileScanningRuleUpdate,
		DeleteContext: fileScanningRuleDelete,
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
			"block_all_file_types": {
				Description:   blockAllFileTypesDesc,
				Type:          schema.TypeBool,
				ConflictsWith: []string{"block_file_types", "block_countries", "block_content_types", "block_threat_types"},
				Optional:      true,
			},
			"block_content_types": {
				Description: blockContentTypesDesc,
				Type:        schema.TypeList,
				MaxItems:    50,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM(common.ContentTypes...),
				},
				Optional: true,
			},
			"block_countries": {
				Description: blockCountriesDesc,
				Type:        schema.TypeList,
				MaxItems:    50,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM(common.Countries...),
				},
				Optional: true,
			},
			"block_file_types": {
				Description: blockFileTypeDesc,
				Type:        schema.TypeList,
				MaxItems:    50,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("ace", "arj", "bat", "cab", "chm", "contact",
						"cpl", "csv", "dll", "doc", "docb", "docm", "docx", "dot", "dotm", "dotx", "exe", "gz", "gzip",
						"hta", "htm", "html", "img", "inetloc", "iqy", "iso", "jar", "jnlp", "js", "lnk", "mam", "mht",
						"mhtml", "msi", "odp", "ods", "odt", "one", "onepkg", "pdf", "php", "plist", "pot", "potm", "potx",
						"ppa", "ppam", "pps", "ppsm", "ppsx", "ppt", "pptm", "pptx", "ps", "ps1", "ps1xml", "ps2", "ps2xml",
						"psc1", "psc2", "pub", "py", "rar", "reg", "rtf", "sh", "shtm", "shtml", "slk", "swf", "vb", "vbe",
						"vbs", "vbscript", "vcard", "vcf", "vcs", "vhd", "vhdx", "wll", "wmv", "xht", "xla", "xlam", "xll",
						"xlm", "xls", "xlsb", "xlsm", "xlsx", "xlt", "xltm", "xltx", "xps", "xxe", "zip"),
				},
				Optional: true,
			},
			"block_threat_types": {
				Description: blockThreatTypesDesc,
				Type:        schema.TypeList,
				MaxItems:    50,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("Abused TLD", "Bitcoin Related", "Blackhole",
						"Botnets", "Brute Forcer", "Chat Server", "CnC", "Compromised", "DDoS Target", "Drop", "DynDNS",
						"EXE Source", "Fake AV", "IP Check", "Keyloggers and Monitoring", "Malware Sites", "Mobile CnC",
						"Mobile Spyware CnC", "Online Gaming", "P2P CnC", "Peer to Peer", "Parking",
						"Phishing and Other Frauds", "Private IP Addresses", "Proxy Avoidance and Anonymizers",
						"Remote Access Service", "Scanner", "Self Signed SSL", "SPAM URLs", "Spyware and Adware",
						"Tor", "Undesirable", "Utility", "VPN"),
				},
				Optional: true,
			},
			"block_unsupported_files": {
				Description: blockUnsupportedFilesDesc,
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"filter_expression": {
				Description: expressionDesc,
				Type:        schema.TypeString,
				Optional:    true,
			},
			"malware": {
				Description:      malwareDesc,
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateStringENUM("DOWNLOAD", "UPLOAD", "ALL"),
			},
			"max_file_size_mb": {
				Description:      maxFileSizeMBDesc,
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: common.ValidateIntRange(1, 4096),
			},
			"priority": {
				Description:      priorityDesc,
				Type:             schema.TypeInt,
				ValidateDiagFunc: common.ValidateIntRange(1, 5000),
				Required:         true,
			},
			"sandbox_file_types": {
				Description: sandboxFileTypesDesc,
				Type:        schema.TypeList,
				MaxItems:    50,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM("cmd", "com", "exe", "msi", "rar"),
				},
				Optional: true,
			},
			"timeout_policy": {
				Description:      timeoutPolicyDesc,
				Type:             schema.TypeString,
				ValidateDiagFunc: common.ValidateStringENUM("PASS", "BLOCK"),
				Optional:         true,
			},
		},
	}
}
