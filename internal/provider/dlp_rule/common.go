package dlp_rule

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
	description = `Data Loss Prevention (DLP) mechanism provides a set of rules aimed at monitoring user behavior and preventing file-based data leaks.
This is achieved by blocking or allowing admin-defined activity (file download or upload) if selected criteria were matched.
The DLP rules include the following parameters:

- Resource (cloud app type, threat, content or country of origin).
- File type.
- Detectors as defined via a separate Proofpoint’s Data Loss Prevention application.

The CASB proxy DLP takes full advantage of Proofpoint’s robust Enterprise DLP solution with its people-centric, scalable, cloud-based architecture.`
	actionDesc = "Enum: `BLOCK`, `LOG`.\n" +
		"This action determines what should be done in case a user performs a user action on files based on this DLP rule."
	applyToOrgDesc           = "Indicates whether this dlp rule applies to the org."
	sourcesDesc              = "Users and groups on which the DLP rule should be applied."
	exemptSources            = "Subgroup of 'sources' on which the DLP rule should not be applied."
	alertLevelDesc           = "ENUM: `CRITICAL`, `HIGH`, `MEDIUM`, `LOW`. Alert severity level for this DLP rule."
	allResourcesDesc         = "Indicates whether to evaluate all resources for scanning."
	allSupportedFileTypeDesc = "Indicates whether all supported file types should get scanned."
	contentTypesDesc         = "A list of content types. " +
		"If a URL is found to be categorized under at least of one of them, file should get scanned. " +
		"ENUM: " + common.ContentTypesDoc
	detectorsDesc = "A list of detectors. When one of the detectors is found during file scan the action will be applied."
	filePartsDesc = "A list of file parts to scan."
	fileTypesDesc = "A list of file types to scan. ENUM: `arj`, `csv`, `doc`, `docm`, `docx`, `dot`, " +
		"`dotm`, `dotx`, `java`, `odp`, `ods`, `odt`, `pdf`, `perl`, `pot`, `potm`, `potx`, `ppa`, `ppam`, `pps`, " +
		"`ppsm`, `ppsx`, `ppt`, `pptm`, `pptx`, `py`, `rtf`, `sh`, `sql`, `txt`, `xla`, `xlam`, `xlm`, `xls`, `xlsb`, " +
		"`xlsm`, `xlsx`, `xlt`, `xltm`, `xltx`, `zip`"
	expressionDesc = `Defines filtering expressions to ensure granularity in DLP rule application.
These expressions consist of the **{Key:Value}** tags according to the internal and external risk factors obtained from the following sources:

	- Proofpoint’s Nexus People Risk Explorer (NPRE).
	- Proofpoint’s Targeted Attack Protection (TAP).
	- CrowdStrike’s Falcon Zero Trust Assessment (ZTA).
	- Configured posture checks.
	- User-defined tags.
	- Auto-generated tags, such as platform type, device type, etc.`
	priorityDesc = "Determines the order in which the DLP rules are evaluated. " +
		"The order is significant since among all the dlp rules that are relevant to a specific user, " +
		"the one with the highest priority (smaller priority value) is the one to determine the dlp enforcement applied to that user."
	resourceCountriesDesc = "A list of countries in which this rule should be applied. Each country should be represented by a Alpha-2 code (ISO-3166). " +
		"Enum: " + common.CountriesDoc
	threatTypesDesc = "A list of threat types. " +
		"If a URL is found to be categorized under at least of one of them, file should get scanned. " +
		"ENUM: `Abused TLD`, `Bitcoin Related`, `Blackhole`, `Botnets`, `Brute Forcer`, `Chat Server`, `CnC`, " +
		"`Compromised`, `DDoS Target`, `Drop`, `DynDNS`, `EXE Source`, `Fake AV`, `IP Check`, " +
		"`Keyloggers and Monitoring`, `Malware Sites`, `Mobile CnC`, `Mobile Spyware CnC`, `Online Gaming`, `P2P CnC`, " +
		"`Peer to Peer`, `Parking`, `Phishing and Other Frauds`, `Private IP Addresses`, " +
		"`Proxy Avoidance and Anonymizers`, `Remote Access Service`, `Scanner`, `Self Signed SSL`, `SPAM URLs`, " +
		"`Spyware and Adware`, `Tor`, `Undesirable`, `Utility`, `VPN`"
	cloudAppsDesc   = "List of [cloud app](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/cloud_app) IDs on which to apply the dlp rule. "
	userActionsDesc = "A list of user actions on files for which to scan the file. ENUM: `DOWNLOAD`, `UPLOAD`, `CREATE`, `EDIT`, `SHARE`, `POST`, `DELETE`, `LIKE`."
)

var excludedKeys = []string{"id"}

func dlpRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	a, err := client.GetDLPRule(ctx, c, id)
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
func dlpRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewDLPRule(d)
	a, err := client.CreateDLPRule(ctx, c, body)
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

func dlpRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewDLPRule(d)
	a, err := client.UpdateDLPRule(ctx, c, id, body)
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

func dlpRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteDLPRule(ctx, c, id)
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
