package trusted_network

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

const (
	description = `The trusted networks feature is a mechanism for auto-disconnecting from Proofpoint NaaS when the device is on a trusted network, such as corporate environment. The moment the device leaves the trusted network, it auto-reconnects to the Proofpoint NaaS.
A user can still forcefully connect the device to the Proofpoint NaaS when on a trusted network, by clicking Connect in the Proofpoint Agent UI.

Proofpoint Agent tries to detect if a trusted network is available, before re-connecting to a network (re-connect occurs if the device was connected to the Proofpoint NaaS and then switched networks).
A trusted network is defined according to one of the following criteria:
- DNS resolution of a hostname: When a specific hostname is resolved to an IP address that is within the defined IP range.
- External IP address of the device: When the device external IP address is within the specified IP address range.`
	applyToOrgDesc        = "Indicates whether this Trusted Network setting applies to the whole org. Note: This attribute overrides `apply_to_entities`."
	applyToEntitiesDesc   = "Entities (users, groups or network elements) to be subjected to the Trusted Network."
	exemptEntitiesDesc    = "Entities (users, groups or network elements) which are exempt from the Trusted Network."
	externalIPConfigDesc  = "Specified IP address range to compare with the device's external IP for the network to be trusted"
	resolvedAddressConfig = "A hostname and specified IP address range in which the hostname should be resolved for the network to be trusted."
)

var excludedKeys = []string{"id", "criteria"}

func trustedNetworkToResource(d *schema.ResourceData, tn *client.TrustedNetwork) (diags diag.Diagnostics) {
	err := client.MapResponseToResource(tn, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	criteriaList := make([]interface{}, len(tn.Criteria))
	for i, cr := range tn.Criteria {
		criteria := make(map[string]interface{})
		criteria["type"] = cr.Type
		if cr.ResolvedAddressConfig != nil {
			rac := make([]interface{}, 1)
			resolvedAddressConfig := make(map[string]interface{})
			resolvedAddressConfig["hostname"] = cr.ResolvedAddressConfig.Hostname
			addressRange := make([]string, len(cr.ResolvedAddressConfig.AddressesRanges))
			for j, address := range cr.ResolvedAddressConfig.AddressesRanges {
				addressRange[j] = address
			}
			resolvedAddressConfig["addresses_ranges"] = addressRange
			rac[0] = resolvedAddressConfig
			criteria["resolved_address_config"] = rac
		}
		if cr.ExternalIpConfig != nil {
			eic := make([]interface{}, 1)
			externalIpConfig := make(map[string]interface{})
			addressRange := make([]string, len(cr.ExternalIpConfig.AddressesRanges))
			for j, address := range cr.ExternalIpConfig.AddressesRanges {
				addressRange[j] = address
			}
			externalIpConfig["addresses_ranges"] = addressRange
			eic[0] = externalIpConfig
			criteria["external_ip_config"] = eic
		}
		criteriaList[i] = criteria
	}
	if err = d.Set("criteria", criteriaList); err != nil {
		return diag.FromErr(err)
	}
	return
}

func trustedNetworkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	tn, err := client.GetTrustedNetwork(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId(tn.ID)
	return trustedNetworkToResource(d, tn)
}
func trustedNetworkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewTrustedNetwork(d)
	tn, err := client.CreateTrustedNetwork(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(tn.ID)
	return trustedNetworkToResource(d, tn)
}

func trustedNetworkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewTrustedNetwork(d)
	tn, err := client.UpdateTrustedNetwork(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(tn.ID)
	return trustedNetworkToResource(d, tn)
}

func trustedNetworkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteTrustedNetwork(ctx, c, id)
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
