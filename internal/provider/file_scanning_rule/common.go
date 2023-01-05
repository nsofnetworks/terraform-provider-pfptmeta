package file_scanning_rule

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
	"log"
	"net/http"
)

const (
	description = `Enhanced file scanning techniques protect organizations and their end users from downloading the following objects:

- Malware.
- Files from undesirable sources.
- Types of files prohibited by organization.`
	applyToOrgDesc = "Indicates whether this file scanning rule applies to the org."
	sourcesDesc    = "Users and groups on which the file scanning rule should be applied."
	exemptSources  = "Subgroup of 'sources' on which the file scanning rule should not be applied."
	cloudAppsDesc  = "List of cloud apps on which to apply the file scanning rule."
	expressionDesc = `Defines filtering expressions to ensure granularity in file scanning rule application.
These expressions consist of the **{Key:Value}** tags according to the internal and external risk factors obtained from the following sources:

	- Proofpoint’s Nexus People Risk Explorer (NPRE).
	- Proofpoint’s Targeted Attack Protection (TAP).
	- CrowdStrike’s Falcon Zero Trust Assessment (ZTA).
	- Configured posture checks.
	- User-defined tags.
	- Auto-generated tags, such as platform type, device type, etc.`
	blockAllFileTypesDesc = "Indicates whether all file types should get blocked."
	blockContentTypesDesc = "A List of content types. If a URL is found to be categorized under at least of one of them, all file types should get blocked. Enum: " + common.ContentTypesDoc
	blockCountriesDesc    = "A list of countries to which all file types should get blocked. Each country should be represented by a Alpha-2 code (ISO-3166). Enum: " + common.CountriesDoc
	blockFileTypeDesc     = "A List of file types to block. Enum: `ace`, `arj`, `bat`, `cab`, `chm`, `contact`, `cpl`, " +
		"`csv`, `dll`, `doc`, `docb`, `docm`, `docx`, `dot`, `dotm`, `dotx`, `exe`, `gz`, `gzip`, `hta`, `htm`, `html`, " +
		"`img`, `inetloc`, `iqy`, `iso`, `jar`, `jnlp`, `js`, `lnk`, `mam`, `mht`, `mhtml`, `msi`, `odp`, `ods`, `odt`, " +
		"`one`, `onepkg`, `pdf`, `php`, `plist`, `pot`, `potm`, `potx`, `ppa`, `ppam`, `pps`, `ppsm`, `ppsx`, `ppt`, " +
		"`pptm`, `pptx`, `ps`, `ps1`, `ps1xml`, `ps2`, `ps2xml`, `psc1`, `psc2`, `pub`, `py`, `rar`, `reg`, `rtf`, `sh`, " +
		"`shtm`, `shtml`, `slk`, `swf`, `vb`, `vbe`, `vbs`, `vbscript`, `vcard`, `vcf`, `vcs`, `vhd`, `vhdx`, `wll`, `wmv`, " +
		"`xht`, `xla`, `xlam`, `xll`, `xlm`, `xls`, `xlsb`, `xlsm`, `xlsx`, `xlt`, `xltm`, `xltx`, `xps`, `xxe`, `zip`"
	blockThreatTypesDesc = "A List of threat types. If a URL is found to be categorized under at least of one of them, all file types should get blocked. " +
		"Enum: `Abused TLD`, `Bitcoin Related`, `Blackhole`, `Botnets`, `Brute Forcer`, `Chat Server`, `CnC`, " +
		"`Compromised`, `DDoS Target`, `Drop`, `DynDNS`, `EXE Source`, `Fake AV`, `IP Check`, `Keyloggers and Monitoring`, " +
		"`Malware Sites`, `Mobile CnC`, `Mobile Spyware CnC`, `Online Gaming`, `P2P CnC`, `Peer to Peer`, `Parking`, " +
		"`Phishing and Other Frauds`, `Private IP Addresses`, `Proxy Avoidance and Anonymizers`, `Remote Access Service`, " +
		"`Scanner`, `Self Signed SSL`, `SPAM URLs`, `Spyware and Adware`, `Tor`, `Undesirable`, `Utility`, `VPN`"
	blockUnsupportedFilesDesc = "Indicates whether unsupported files should get blocked. It includes unsupported file types, " +
		"files that are too big to be scanned and possibly other unsupported conditions."
	malwareDesc       = "Indicates whether malware should be scanned for upload and or download of files. Enum: `DOWNLOAD`, `UPLOAD`, `ALL`"
	maxFileSizeMBDesc = "The maximal size of a file in MB to scan. Any file larger than this threshold will get blocked. " +
		"If not specified, no limit on maximal file size is enforced."
	priorityDesc = "Determines the order in which the file scanning rules are evaluated." +
		" The order is significant since among all the file scanning rules that are relevant to a specific user, " +
		"the one with the highest priority (smaller priority value) is the one to determine the file scanning enforcement applied to that user."
	sandboxFileTypesDesc = "A List of file types to sandbox. Enum: `cmd`, `com`, `exe`, `msi`, `rar`"
	timeoutPolicyDesc    = "Whether to prevent or allow execution of the current file handling rule if it runs over the time limit (1 minute). Enum: `PASS`, `BLOCK`"
)

var excludedKeys = []string{"id"}

func fileScanningRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	a, err := client.GetFileScanningRule(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing dlp rule %s because it's gone", id)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId(a.ID)
	err = client.MapResponseToResource(a, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}
func fileScanningRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewFileScanningRule(d)
	a, err := client.CreateFileScanningRule(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(a.ID)
	err = client.MapResponseToResource(a, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func fileScanningRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewFileScanningRule(d)
	a, err := client.UpdateFileScanningRule(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(a.ID)
	err = client.MapResponseToResource(a, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func fileScanningRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteFileScanningRule(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return
}
