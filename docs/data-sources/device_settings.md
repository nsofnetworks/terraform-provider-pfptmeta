---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Data Source pfptmeta_device_settings - terraform-provider-pfptmeta"
subcategory: "Device Management"
description: |-
  The pfptmeta_device_settings resource is a tool with which the administrator can configure user devices.
  The settings are related to authentication, access and security that can be defined for a specific device, a user, a group of users, or the entire organization.
---

# Data Source (pfptmeta_device_settings)

The `pfptmeta_device_settings` resource is a tool with which the administrator can configure user devices.
The settings are related to authentication, access and security that can be defined for a specific device, a user, a group of users, or the entire organization.

## Example Usage

```terraform
data "pfptmeta_device_settings" "settings" {
  id = "ds-123abc"
}

output "settings" {
  value = data.pfptmeta_device_settings.settings
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Read-Only

- `apply_on_org` (Boolean) Indicates whether this device setting applies to the entire org. Note: this attribute overrides `apply_to_entities`.
- `apply_to_entities` (List of String) Entities (users, groups or network elements) that the device settings will be applied to.
- `auto_fqdn_domain_names` (List of String) Auto-generated FQDNs of devices concatenated by {hostname}.{domain_name} (hostname is reported by the agent). Use `[""]` to utilize the reported hostname, or add a domain name to be concatenated with the hostname, or omit to disable it.
- `description` (String)
- `direct_sso` (String) User authentication is enforced via the selected IdP. The user will be automatically redirected to the IdP login page for authentication. Uses the Identity Provider ID.
- `enabled` (Boolean)
- `id` (String) The ID of this resource.
- `name` (String)
- `overlay_mfa_refresh_period` (Number) User auth-token lifetime in minutes. During auth-token lifetime, users can (re)connect without entering login credentials. Must be >= 10.
- `overlay_mfa_required` (Boolean) Defines whether users need to authenticate with their login credentials when they connect. If not required, the authentication is done only with the user's client certificate.
- `protocol_selection_lifetime` (String) Integer wrapped as string. A time period (in minutes) after which the Proofpoint Agent attempts to reconnect using IPsec after previous automatic switchover to TLS.
- `proxy_always_on` (Boolean) Controls the Web Security always-on enforcement on end-user devices.
- `search_domains` (List of String) Domain search list. These domains are used by the device resolver to create a Fully Qualified Domain Name (FQDN) from a relative name. The resolver tries resolving the search domains in the order they are listed. If all resolutions fail, it attempts to resolve the original query name.
- `session_lifetime` (Number) Specifies the number of minutes allowed for a user session, since the last user authentication. The session is terminated when session lifetime expires. Must be >= 1.
- `session_lifetime_grace` (String) Integer wrapped as string. Specifies the number of minutes for a user to get notified before the session is about to expire due to Session Lifetime. Must be between 0 to 60.
- `tunnel_mode` (String) Specifies the tunnel operation mode:
	- **split** - Internet traffic is not tunneled, and only traffic to private (mapped) resources is routed through Proofpoint NaaS.
	- **full**- All traffic is tunneled and routed through Proofpoint NaaS.
- `vpn_login_browser` (String) Forces login to VPN via agent/external_browser or lets the user to decide. Enum: `AGENT`, `EXTERNAL`, `USER_DEFINED`.
- `ztna_always_on` (Boolean) Defines whether to enforce a persistent connection of end-user device to the ZTNA network.
