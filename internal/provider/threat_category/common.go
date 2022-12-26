package threat_category

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
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
		"Enum: " + common.CountriesDoc
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
