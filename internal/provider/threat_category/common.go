package threat_category

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

var excludedKeys = []string{"id"}

const (
	description = "The administrators can enhance filtering capabilities by adding web threat categories to the access rules. " +
		"By creating objects for triggering a violation based on threat categories, " +
		"the administrator can protect their users from hostile activities by reducing malware infections, " +
		"neutralizing phishing attempts, preventing botnet takeovers and so on. " +
		"Proofpoint Meta provides three default threat classes (permissive, moderate and strict), " +
		"each with its own list of preset web threats."
	confidenceLevelDesc = "Provides the accuracy degree for recognizing the selected traffic type as threat, " +
		"as defined by the security engines. " +
		"This can be used to reduce potential false-positives or fine-tune the system to suit better for the company’s specific needs. " +
		"By enabling this feature, the administrator defines a tolerance threshold (low, medium or high). " +
		"When the threshold is crossed, a rule violation is triggered. " +
		"ENUM: `LOW`, `MEDIUM`, `HIGH`."
	riskLevelDesc = "Indicates the risk level that the security engines have for any particular site. " +
		"By enabling this feature, the administrator sets a tolerance threshold (low, medium or high). " +
		"When the threshold is crossed, a rule violation is triggered. " +
		"ENUM: `LOW`, `MEDIUM`, `HIGH`."
	countriesDesc = "A list of countries to which access should be restricted. Each country should be represented by a Alpha-2 code (ISO-3166). " +
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
	typesDesc = "A list of predefined threat types to protect against. " +
		"Enum:`Abused TLD`,`Bitcoin Related`,`Blackhole`,`Botnets`,`Brute Forcer`,`Chat Server`,`CnC`," +
		"`Compromised`,`DDoS Target`,`Drop`,`DynDNS`,`EXE Source`,`Fake AV`,`IP Check`,`Keyloggers and Monitoring`," +
		"`Malware Sites`,`Mobile CnC`,`Mobile Spyware CnC`,`Online Gaming`,`P2P CnC`,`Peer to Peer`,`Parking`," +
		"`Phishing and Other Frauds`,`Private IP Addresses`,`Proxy Avoidance and Anonymizers`,`Remote Access Service`," +
		"`Scanner`,`Self Signed SSL`,`SPAM URLs`,`Spyware and Adware`,`Tor`,`Undesirable`,`Utility`,`VPN`"
)

func threatCategoryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	tc, err := client.GetThreatCategory(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing threat category %s because it's gone", id)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId(tc.ID)
	err = client.MapResponseToResource(tc, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}
func threatCategoryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewThreatCategory(d)
	tc, err := client.CreateThreatCategory(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(tc.ID)
	err = client.MapResponseToResource(tc, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func threatCategoryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewThreatCategory(d)
	tc, err := client.UpdateThreatCategory(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(tc.ID)
	err = client.MapResponseToResource(tc, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func threatCategoryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteThreatCategory(ctx, c, id)
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
