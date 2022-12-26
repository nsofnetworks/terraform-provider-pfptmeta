package threat_category

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: threatCategoryCreate,
		ReadContext:   threatCategoryRead,
		UpdateContext: threatCategoryUpdate,
		DeleteContext: threatCategoryDelete,
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
			"confidence_level": {
				Description:      confidenceLevelDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateStringENUM("LOW", "MEDIUM", "HIGH"),
			},
			"risk_level": {
				Description:      riskLevelDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateStringENUM("LOW", "MEDIUM", "HIGH"),
			},
			"types": {
				Description: typesDesc,
				Type:        schema.TypeList,
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
				Optional:     true,
				AtLeastOneOf: []string{"types", "countries"},
			},
			"countries": {
				Description: countriesDesc,
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateStringENUM(common.Countries...),
				},
				Optional:     true,
				AtLeastOneOf: []string{"types", "countries"},
			},
		},
	}
}
