package device_settings

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
	"strconv"
)

var excludedKeys = []string{"id", "protocol_selection_lifetime", "session_lifetime_grace"}

const (
	description = "The `pfptmeta_device_settings` resource is a tool with which the administrator can configure user devices.\n" +
		"The settings are related to authentication, access and security that can be defined for a specific device, a user, a group of users, or the entire organization."
	applyOnOrgDesc                = "Indicates whether this device setting applies to the entire org. Note: this attribute overrides `apply_to_entities`."
	applyToEntitiesDesc           = "Entities (users, groups or network elements) that the device settings will be applied to."
	autoFqdnDomainNamesDesc       = "Auto-generated FQDNs of devices concatenated by {hostname}.{domain_name} (hostname is reported by the agent). Use `[\"\"]` to utilize the reported hostname, or add a domain name to be concatenated with the hostname, or omit to disable it."
	directSsoDesc                 = "User authentication is enforced via the selected IdP. The user will be automatically redirected to the IdP login page for authentication. Uses the Identity Provider ID."
	overlayMFARefreshPeriodDesc   = "User auth-token lifetime in minutes. During auth-token lifetime, users can (re)connect without entering login credentials. Must be >= 10."
	overlayMFARequiredDesc        = "Defines whether users need to authenticate with their login credentials when they connect. If not required, the authentication is done only with the user's client certificate."
	protocolSelectionLifeTimeDesc = "Integer wrapped as string. A time period (in minutes) after which the Proofpoint Agent attempts to reconnect using IPsec after previous automatic switchover to TLS."
	searchDomainsDesc             = "Domain search list. These domains are used by the device resolver to create a Fully Qualified Domain Name (FQDN) from a relative name. The resolver tries resolving the search domains in the order they are listed. If all resolutions fail, it attempts to resolve the original query name."
	proxyAlwaysOnDesc             = "Controls the Web Security always-on enforcement on end-user devices."
	sessionLifeTimeDesc           = "Specifies the number of minutes allowed for a user session, since the last user authentication. The session is terminated when session lifetime expires. Must be >= 1."
	sessionLifeTimeGraceDesc      = "Integer wrapped as string. Specifies the number of minutes for a user to get notified before the session is about to expire due to Session Lifetime. Must be between 0 to 60."
	tunnelModeDesc                = "Specifies the tunnel operation mode:\n" +
		"	- **split** - Internet traffic is not tunneled, and only traffic to private (mapped) resources is routed through Proofpoint NaaS.\n" +
		"	- **full**- All traffic is tunneled and routed through Proofpoint NaaS."
	vpnLoginBrowserDesc = "Forces login to VPN via agent/external_browser or lets the user to decide. Enum: `AGENT`, `EXTERNAL`, `USER_DEFINED`."
	ztnaAlwaysOnDesc    = "Defines whether to enforce a persistent connection of end-user device to the ZTNA network."
)

func deviceSettingsToResource(d *schema.ResourceData, ds *client.DeviceSettings) (diags diag.Diagnostics) {
	d.SetId(ds.ID)
	err := client.MapResponseToResource(ds, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	if ds.ProtocolSelectionLifetime == nil {
		err = d.Set("protocol_selection_lifetime", nil)
	} else {
		err = d.Set("protocol_selection_lifetime", strconv.Itoa(*ds.ProtocolSelectionLifetime))
	}
	if err != nil {
		return diag.FromErr(err)
	}
	if ds.SessionLifetimeGrace == nil {
		err = d.Set("session_lifetime_grace", nil)
	} else {
		err = d.Set("session_lifetime_grace", strconv.Itoa(*ds.SessionLifetimeGrace))
	}
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func deviceSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	ds, err := client.GetDeviceSettings(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	return deviceSettingsToResource(d, ds)
}
func deviceSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewDeviceSettings(d)
	ds, err := client.CreateDeviceSettings(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return deviceSettingsToResource(d, ds)
}

func deviceSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewDeviceSettings(d)
	ds, err := client.UpdateDeviceSettings(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return deviceSettingsToResource(d, ds)
}

func deviceSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteDeviceSettings(ctx, c, id)
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
