package scan_rule

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
	description = `Scan mechanism provides a set of rules aimed at monitoring user behavior and preventing file-based data leaks.
This is achieved by blocking or allowing admin-defined activity if selected criteria were matched.
The Scan rules include the following parameters:

- Resource.
- Malware.
- File type.
- Detectors as defined via a separate Proofpoint’s Data Loss Prevention application.

The CASB proxy Scan takes full advantage of Proofpoint’s robust Enterprise Scan solution with its people-centric, scalable, cloud-based architecture.`
	priorityDesc = "Determines the order in which the Scan rules are evaluated. " +
		"The order is significant since among all the scan rules that are relevant to a specific user, " +
		"the one with the highest priority (smaller priority value) is the one to determine the scan enforcement applied to that user."
	sourcesDesc       = "Users and groups on which the Scan rule should be applied."
	exemptSourcesDesc = "Subgroup of 'sources' on which the Scan rule should not be applied."
	expressionDesc    = `Defines filtering expressions to ensure granularity in Scan rule application.
These expressions consist of the **{Key:Value}** tags according to the internal and external risk factors obtained from the following sources:

	- Proofpoint’s Nexus People Risk Explorer (NPRE).
	- Proofpoint’s Targeted Attack Protection (TAP).
	- CrowdStrike’s Falcon Zero Trust Assessment (ZTA).
	- Configured posture checks.
	- User-defined tags.
	- Auto-generated tags, such as platform type, device type, etc.`
	applyToOrgDesc         = "Indicates whether this scan rule applies to the org."
	contentCategoriesDesc  = "List of [content category](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/content_category) IDs which the Scan rule should process."
	threatCategoriesDesc   = "List of [threat category](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/threat_category) IDs the Scan rule will protect against"
	cloudAppsDesc          = "List of [cloud app](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/cloud_app) IDs on which to apply the scan rule. "
	cloudAppRiskGroupsDesc = "ENUM: `imminent_targets`, `latent_targets`, `major_targets`, `soft_targets`, " +
		"`very_attacked_apps`, `very_privileged_apps`, `very_vulnerable_apps`\n" +
		"List of cloud app NCRE based risk groups which the scan rule should process."
	catalogAppCategoriesDesc = "ENUM: `Instant Messaging`, `eCommerce`, `Content Management`, `Software Development`, `Project Management`, " +
		"`Marketing`, `CRM`, `Telecommunications`, `Social and Communication`, `Productivity`, `Collaboration`, " +
		"`Business and Finance`, `Utilities`, `IT Service Management`, `Social Networking`, `Office Document and Productivity`, " +
		"`Cloud File Sharing`, `Web Meetings`, `Identity and Access Management`, `IT Services and Hosting`, `Webmail`, " +
		"`Website Builder`, `Human Capital Management`, `Sales and CRM`, `E-commerce and Accounting`, `Streaming Media`, " +
		"`Cloud Storage`, `Operations Management`, `Online Meeting`, `Supply Chain`, `Security and Compliance`, " +
		"`Entertainment and Lifestyle`, `System and Network`, `Retail and Consumer Services`, `Health and Benefits`, " +
		"`Data and Analytics`, `Education and References`, `Personal instant messaging`, `Legal`, `Other`, `Hosting Services`, " +
		"`News and Media`, `Sales`, `Enterprise Resource Planning`, `Advertising`, `Travel and Transportation`, " +
		"`Property Management`, `Government Services`, `Games`, `Code Hosting`.\n" +
		"List of catalog app categories that the Scan rule must process."
	networksDesc  = "List of source [IP network](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/ip_network) IDs the Scan rule applies on"
	countriesDesc = "A list of source countries in which this rule should be applied." +
		"Each country should be represented by a Alpha-2 code (ISO-3166)." +
		"Enum: " + common.CountriesDoc
	userAgentsDesc = "ENUM: `Chrome`, `Safari`, `Edge`, `Firefox`, `Opera`, `IE`, `Electron`, `Outlook`, `Excel`, `PowerPoint\n`" +
		"A List of user agents on which the rule applies.\n" +
		"Meaning, in order for the rule to be evaluated, the user agent that was used to make the request must be on that list."
	userActionsDesc          = "A list of user actions on files for which to scan the file. ENUM: `DOWNLOAD`, `UPLOAD`, `CREATE`, `EDIT`, `SHARE`, `POST`, `DELETE`, `LIKE`."
	allSupportedFileTypeDesc = "Indicates whether all supported file types should get scanned."
	fileTypesDesc            = "A list of file types to scan. ENUM: `7z`, `ace`, `arj`, `bat`, `cab`, `chm`, `contact`, `cpl`, `csv`, `dll`, `dmg`, `doc`, `docb`, `docm`, `docx`, `dot`, " +
		"`dotm`, `dotx`, `exe`, `form`, `gz`, `gzip`, `hta`, `htm`, `html`, `img`, `inetloc`, `iqy`, `iso`, `jar`, `java`, `jnlp`, `js`, `lnk`, `mam`, `mht`, `mhtml`, `msi`, `odp`, `ods`, `odt`, " +
		"`one`, `onepkg`, `pdf`, `perl`, `php`, `pkg`, `plist`, `pot`, `potm`, `potx`, `ppa`, `ppam`, `pps`, `ppsm`, `ppsx`, `ppt`, `pptm`, `pptx`, `ps`, `ps1`, `ps1xml`, `ps2`, `ps2xml`, `psc1`, " +
		"`psc2`, `pub`, `py`, `rar`, `reg`, `rtf`, `sh`, `shtm`, `shtml`, `slk`, `sql`, `swf`, `txt`, `vb`, `vbe`, `vbs`, `vbscript`, `vcard`, `vcf`, `vcs`, `vhd`, `vhdx`, `wll`, `wmv`, `xht`, `xla`, " +
		"`xlam`, `xll`, `xlm`, `xls`, `xlsb`, `xlsm`, `xlsx`, `xlt`, `xltm`, `xltx`, `xps`, `xxe`, `zip`"
	maxFileSizeMBDesc = "The maximal size of a file in MB to scan. Any file larger than this threshold will get processed. " +
		"If not specified, no limit on maximal file size is enforced."
	passwordProtectedFilesDesc = "Indicates whether password protected files should get processed."
	dlpDesc                    = "Enable dlp file scan according to the list of detectors."
	detectorsDesc              = "A list of detectors. When one of the detectors is found during file scan the action will be applied."
	malwareDesc                = "Indicates whether malware should be scanned for upload and or download of files according to the defined user actions."
	sandboxDesc                = "Indicates whether files should be sandboxed. Only relevant if malware is enabled."
	antivirusDesc              = "Indicates whether files should be scanned by antivirus. Only relevant if malware is enabled."
	actionDesc                 = "Enum: `BLOCK`, `LOG`.\n" +
		"This action determines what should be done in case a user performs a user action on files based on this Scan rule."
	AccessIdsDesc = "Devices on which the Scan rule should be applied"
)

var excludedKeys = []string{"id"}

func scanRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	a, err := client.GetScanRule(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing scan rule %s because it's gone", id)
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
func scanRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewScanRule(d)
	a, err := client.CreateScanRule(ctx, c, body)
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

func scanRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewScanRule(d)
	a, err := client.UpdateScanRule(ctx, c, id, body)
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

func scanRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteScanRule(ctx, c, id)
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
