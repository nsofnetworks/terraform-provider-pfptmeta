package url_filtering_rule

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

const (
	description = `The Proofpoint Web Security solution protects against web-based security threats by defining URL filtering rules
which include various content and threat categories, as well as cloud-based applications and tenant restrictions.
With these measures, you can enforce company security policies and filter malicious internet traffic in real time.`
	actionDesc = "Enum: `ISOLATION`, `BLOCK`, `LOG`, `RESTRICT`, `WARN`.\n" +
		"This action determines what must be done according to this URL filtering rule if a user tries to reach a restricted URL."
	applyToOrgDesc               = "indicates whether this URL filtering rule applies to the org."
	sourcesDesc                  = "Users and groups on which the URL filtering rule should be applied."
	exemptSources                = "Subgroup of 'sources' on which the URL filtering rule should not be applied."
	advancedThreatProtectionDesc = "Enables the first-rate security engine based on up-to-date web threat intelligence gathered from two decades of protecting the world's largest organizations from email-borne attacks."
	catalogAppCategories         = "ENUM: `Instant Messaging`, `eCommerce`, `Content Management`, `Software Development`, `Project Management`, " +
		"`Marketing`, `CRM`, `Telecommunications`, `Social and Communication`, `Productivity`, `Collaboration`, " +
		"`Business and Finance`, `Utilities`, `IT Service Management`, `Social Networking`, `Office Document and Productivity`, " +
		"`Cloud File Sharing`, `Web Meetings`, `Identity and Access Management`, `IT Services and Hosting`, `Webmail`, " +
		"`Website Builder`, `Human Capital Management`, `Sales and CRM`, `E-commerce and Accounting`, `Streaming Media`, " +
		"`Cloud Storage`, `Operations Management`, `Online Meeting`, `Supply Chain`, `Security and Compliance`, " +
		"`Entertainment and Lifestyle`, `System and Network`, `Retail and Consumer Services`, `Health and Benefits`, " +
		"`Data and Analytics`, `Education and References`, `Personal instant messaging`, `Legal`, `Other`, `Hosting Services`, " +
		"`News and Media`, `Sales`, `Enterprise Resource Planning`, `Advertising`, `Travel and Transportation`, " +
		"`Property Management`, `Government Services`, `Games`, `Code Hosting`.\n" +
		"List of catalog app categories that the URL filtering rule must restrict."
	catalogAppRiskDesc = "Risk threshold to be used to restrict all catalog apps which has that risk or higher."
	cloudAppsDesc      = "List of [cloud app](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/cloud_app) IDs which the URL filtering rule should restrict. "
	countriesDesc      = "A list of countries in which this rule should be applied. Each country should be represented by a Alpha-2 code (ISO-3166). " +
		"Enum: `AD`,`AE`,`AF`,`AG`,`AI`,`AL`,`AM`,`AO`,`AQ`,`AR`,`AS`,`AT`,`AU`,`AW`,`AX`,`AZ`,`BA`,`BB`,`BD`,`BE`,`BF`," +
		"`BG`,`BH`,`BI`,`BJ`,`BL`,`BM`,`BN`,`BO`,`BQ`,`BR`,`BS`,`BT`,`BV`,`BW`,`BY`,`BZ`,`CA`,`CC`,`CD`,`CF`,`CG`,`CH`," +
		"`CI`,`CK`,`CL`,`CM`,`CN`,`CO`,`CR`,`CU`,`CV`,`CW`,`CX`,`CY`,`CZ`,`DE`,`DJ`,`DK`,`DM`,`DO`,`DZ`,`EC`,`EE`,`EG`," +
		"`EH`,`ER`,`ES`,`ET`,`FI`,`FJ`,`FK`,`FM`,`FO`,`FR`,`GA`,`GB`,`GD`,`GE`,`GF`,`GG`,`GH`,`GI`,`GL`,`GM`,`GN`,`GP`," +
		"`GQ`,`GR`,`GS`,`GT`,`GU`,`GW`,`GY`,`HK`,`HM`,`HN`,`HR`,`HT`,`HU`,`ID`,`IE`,`IL`,`IM`,`IN`,`IO`,`IQ`,`IR`,`IS`," +
		"`IT`,`JE`,`JM`,`JO`,`JP`,`KE`,`KG`,`KH`,`KI`,`KM`,`KN`,`KP`,`KR`,`KW`,`KY`,`KZ`,`LA`,`LB`,`LC`,`LI`,`LK`,`LR`," +
		"`LS`,`LT`,`LU`,`LV`,`LY`,`MA`,`MC`,`MD`,`ME`,`MF`,`MG`,`MH`,`MK`,`ML`,`MM`,`MN`,`MO`,`MP`,`MQ`,`MR`,`MS`,`MT`," +
		"`MU`,`MV`,`MW`,`MX`,`MY`,`MZ`,`NA`,`NC`,`NE`,`NF`,`NG`,`NI`,`NL`,`NO`,`NP`,`NR`,`NU`,`NZ`,`OM`,`PA`,`PE`,`PF`," +
		"`PG`,`PH`,`PK`,`PL`,`PM`,`PN`,`PR`,`PS`,`PT`,`PW`,`PY`,`QA`,`RE`,`RO`,`RS`,`RU`,`RW`,`SA`,`SB`,`SC`,`SD`,`SE`," +
		"`SG`,`SH`,`SI`,`SJ`,`SK`,`SL`,`SM`,`SN`,`SO`,`SR`,`SS`,`ST`,`SV`,`SX`,`SY`,`SZ`,`TC`,`TD`,`TF`,`TG`,`TH`,`TJ`," +
		"`TK`,`TL`,`TM`,`TN`,`TO`,`TR`,`TT`,`TV`,`TW`,`TZ`,`UA`,`UG`,`UM`,`US`,`UY`,`UZ`,`VA`,`VC`,`VE`,`VG`,`VI`,`VN`," +
		"`VU`,`WF`,`WS`,`YE`,`YT`,`ZA`,`ZM`,`ZW`"
	expiresAtDesc = "Defines the rule expiration time. " +
		"This can be useful when creating exceptions for users who need them for a limited period of time as an alternative for full disconnection from the proxy. " +
		"When no value is given the URL filtering rule will never expire. Takes `RFC3339` (`2006-01-02T15:04:05Z`) date format."
	expressionDesc = `Defines filtering expressions to ensure granularity in URL filtering rule application.
These expressions consist of the **{Key:Value}** tags according to the internal and external risk factors obtained from the following sources:

	- Proofpoint’s Nexus People Risk Explorer (NPRE).
	- Proofpoint’s Targeted Attack Protection (TAP).
	- CrowdStrike’s Falcon Zero Trust Assessment (ZTA).
	- Configured posture checks.
	- User-defined tags.
	- Auto-generated tags, such as platform type, device type, etc.
`
	contentCategoriesDesc = "List of [content category](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/content_category) IDs which the URL filtering rule should restrict."
	networkDesc           = "List of source [IP network](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/ip_network) IDs the URL filtering rule applies on"
	priorityDesc          = "Determines the order in which the URL-filtering rules are evaluated. " +
		"The order is significant since the first URL-filtering rule that finds a URL restricted is the one to determine which action to execute. " +
		"Lower priority value means the URL-filtering rule will be evaluated earlier."
	scheduleDesc          = "List of [time frame](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/time_frame) IDs during which the URL filtering rule will be enforced"
	tenantRestrictionDesc = "[Tenant restrictions](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/tenant_restriction) for this rule. " +
		"Only the `RESTRICT` action is allowed when this option is set."
	threatCategoriesDesc = "List of [threat category](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/threat_category) IDs the URL filtering rule will protect against"
	warnTtlDesc          = "Time in minutes during which the warning page is not shown again after user proceeds to URL"
)

var excludedKeys = []string{"id"}

func urlFilteringRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	a, err := client.GetUrlFilteringRule(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing url filtering rule %s because it's gone", id)
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
func urlFilteringRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewUrlFilteringRule(d)
	a, err := client.CreateUrlFilteringRule(ctx, c, body)
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

func urlFilteringRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewUrlFilteringRule(d)
	a, err := client.UpdateUrlFilteringRule(ctx, c, id, body)
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

func urlFilteringRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteUrlFilteringRule(ctx, c, id)
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
