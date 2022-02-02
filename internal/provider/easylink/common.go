package easylink

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

const (
	description = `
Proofpoint provides two access modes for end users:

- Agent-based, allows users access to resources natively, after establishing a VPN connection and authenticating using the Proofpoint Agent client.

- Clientless, or MetaConnect (MC), allows users to access resources via supported browsers, without any client installation.

MetaConnect eliminates the need to install an agent and establish a VPN connection. Users access internal applications from a dedicated web (MetaConnect) portal. After authentication, the users see the list of applications that they are allowed to access. Alternatively, users can use a static FQDN to access their private apps directly. MetaConnect can be used to access web applications (HTTP or HTTPS), or servers via RDP, SSH or VNC.

There are several use cases for using MetaConnect access mode. It is well suited for instances when an agent cannot be used, such as with external contractors or personal/BYO devices.

Clientless applications are defined using EasyLinks. An EasyLink defines the application (internal host, protocol and port), the users assigned to this application, the URL type, etc.

See [here](https://help.metanetworks.com/knowledgebase/easylinks/) for more details.`
	domainNameDesc = "FQDN or IPv4 of the application defined by the EasyLink."
	viewersDesc    = "User or group IDs that will be granted access to the application defined by the EasyLink."
	accessFQDNDesc = "External FQDN to be associated with the current EasyLink, required when `access_type` is set to `redirect` or `native`."
	accessTypeDesc = `When creating an Easylink, you need to select the appropriate access (URL) type.

	- **meta** – Use Proofpoint-generated URL as the application entry point and throughout the browsing session.

	- **redirect** – Use your own (vanity) FQDN as the application access entry point, and redirect to Proofpoint-generated URL for the rest of browsing session.

	- **native** – Use your own (vanity) FQDN as the application entry point and throughout the browsing session. Only web-based apps (HTTP or HTTPS) are supported.`
	mappedElementIDDesc = "Hosting resource for Mapped Subnet or Mapped Service network elements if the host is to reside permanently within this resource. This field is required when the host is an IPv4 address."
	certificateIDDesc   = "Required when `access_type` is set to `redirect` or `native`. MetaConnect provides a secure connection with HTTPS. For the end-user browser to trust the domain, you must generate an SSL certificate for the external application FQDN using the `pfptmeta_certificate` resource."
	auditDesc           = "When enabled, all web traffic is logged to the MetaConnect Web log. Logging is only applicable when `protocol` is either `http` or `https`."
	enableSNIDesc       = "Defines whether to enable SNI or not. The SNI can be enabled only when `protocol` is set to `https`."
	portDesc            = "The port of the application defined by the EasyLink."
	protocolDesc        = "The protocol of the application defined by the EasyLink. ENUM: `ssh`, `rdp`, `vnc`, `http`, `https`."
	rootPathDesc        = "The root path of the application defined by the EasyLink, when `protocol` is `http` or `https`."

	proxyDesc                    = "Additional proxy configuration, available only when `protocol` is set to `http` or `https`."
	proxyEnterpriseAccess        = "When enabled, it resets the session on source IP change to minimize latency if the new source IP has enterprise access to the EasyLink. Allowed only for default ports (80, 443) and when `access_type` is set to `redirect`."
	proxyHostsDesc               = "Additional hosts to be routed to the EasyLink."
	proxyHostHeaderDesc          = "An overwrite to the HTTP host header. It is set to the value of `access_fqdn` when `access_type` is set to `native` and not allowed."
	proxyRewriteContentTypesDesc = "Response content types to be rewritten. ENUM: `html`, `json`, `javascript`, `text`. It is required when `rewrite_hosts` or `rewrite_http` are configured."
	proxyRewriteHosts            = "Defines whether to rewrite hosts in the proxy response to the EasyLink host or not. Rewrites in responses with content type specified in `rewrite_content_types`."
	proxyRewriteHostsClient      = "Selects whether to overwrite hosts in all browser client requests or not."
	proxyRewriteHttpDesc         = "Defines whether to overwrite all `http://` links in proxy response to `https://` or not. Rewrites in responses with content type specified in `rewrite_content_types`."
	proxySharedCookies           = "Selects whether to share cookies between EasyLinks in the same region."

	rdpDesc                    = "Additional RDP configuration, available only when `protocol` is set to `rdp`."
	rdpRemoteAppDesc           = "The remote application to start on the remote desktop. If supported by your remote desktop server, only this application will be visible to the user."
	rdpRemoteAppCmdArgsDesc    = "The command-line arguments, if any, for the remote application."
	rdpRemoteAppWorkDirDesc    = "The working directory, if any, for the remote application."
	rdpRemoteAppSecurityDesc   = "Dictates how data is encrypted and what type of authentication is performed. ENUM: `nla`, `rdp`."
	rdpSeverKeyboardLayoutDesc = "Server-supported keyboard layout. Enum: `english-us`, `german`, `french`, `swiss-french`, `italian`, `japanese`, `swedish`, `unicode`."
)

var excludedKeys = []string{"id", "proxy", "rdp"}

func easyLinkToResource(d *schema.ResourceData, e *client.EasyLink) (diags diag.Diagnostics) {
	d.SetId(e.ID)
	err := client.MapResponseToResource(e, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	if e.Proxy != nil {
		proxy := []map[string]interface{}{
			{
				"enterprise_access":     e.Proxy.EnterpriseAccess,
				"hosts":                 e.Proxy.Hosts,
				"http_host_header":      e.Proxy.HttpHostHeader,
				"rewrite_content_types": e.Proxy.RewriteContentTypes,
				"rewrite_hosts":         e.Proxy.RewriteHosts,
				"rewrite_hosts_client":  e.Proxy.RewriteHostsClient,
				"rewrite_http":          e.Proxy.RewriteHttp,
				"shared_cookies":        e.Proxy.SharedCookies,
			},
		}
		err = d.Set("proxy", proxy)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if e.Rdp != nil {
		rdp := []map[string]interface{}{
			{
				"remote_app":             e.Rdp.RemoteApp,
				"remote_app_cmd_args":    e.Rdp.RemoteAppCmdArgs,
				"remote_app_work_dir":    e.Rdp.RemoteAppWorkDir,
				"security":               e.Rdp.Security,
				"server_keyboard_layout": e.Rdp.ServerKeyboardLayout,
			},
		}
		err = d.Set("rdp", rdp)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		err = d.Set("rdp", nil)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}

func easyLinkRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	e, err := client.GetEasyLink(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	return easyLinkToResource(d, e)
}

func easyLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewEasyLink(d)
	e, err := client.CreateEasyLink(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	diags := easyLinkToResource(d, e)
	if diags.HasError() {
		return diags
	}
	proxy := client.NewProxy(d)
	if proxy != nil {
		err = client.UpdateEasylinkProxy(ctx, c, e.ID, proxy)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	rdp := client.NewRdp(d)
	if rdp != nil {
		err = client.UpdateEasylinkRdp(ctx, c, e.ID, rdp)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return easyLinkRead(ctx, d, c)
}

func easyLinkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewEasyLink(d)

	// Setting protocol to empty string so that it will be omitted from the body as PATCH does not support protocol
	body.Protocol = ""
	e, err := client.UpdateEasyLink(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(e.ID)

	if d.HasChange("proxy") && len(d.Get("proxy").([]interface{})) == 1 {
		proxy := client.NewProxy(d)
		// In case access_type is native - the value is forced to be the value of access_fqdn and cannot be modified.
		// Setting it to empty string will omit it from the request should be removed when NSOF-5878 is resolved
		if e.AccessType == "native" {
			proxy.HttpHostHeader = ""
		}
		err = client.UpdateEasylinkProxy(ctx, c, e.ID, proxy)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("rdp") && len(d.Get("rdp").([]interface{})) == 1 {
		rdp := client.NewRdp(d)
		err = client.UpdateEasylinkRdp(ctx, c, e.ID, rdp)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return easyLinkRead(ctx, d, c)
}

func easyLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteEasyLink(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	return
}
