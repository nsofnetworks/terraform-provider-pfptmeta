package easylink

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,
		ReadContext: easyLinkRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "el"),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Description: domainNameDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"viewers": {
				Description: viewersDesc,
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
			},
			"access_fqdn": {
				Description: accessFQDNDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"access_type": {
				Description: accessTypeDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"mapped_element_id": {
				Description: mappedElementIDDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"certificate_id": {
				Description: certificateIDDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"audit": {
				Description: auditDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"enable_sni": {
				Description: enableSNIDesc,
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"port": {
				Description: portDesc,
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"protocol": {
				Description: protocolDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"root_path": {
				Description: rootPathDesc,
				Type:        schema.TypeString,
				Computed:    true,
			},
			"proxy": {
				Description: proxyDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enterprise_access": {
							Description: proxyEnterpriseAccess,
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"hosts": {
							Description: proxyHostsDesc,
							Computed:    true,
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"http_host_header": {
							Description: proxyHostHeaderDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"rewrite_content_types": {
							Description: proxyRewriteContentTypesDesc,
							Computed:    true,
							Type:        schema.TypeList,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"rewrite_hosts": {
							Description: proxyRewriteHosts,
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"rewrite_hosts_client": {
							Description: proxyRewriteHostsClient,
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"rewrite_http": {
							Description: proxyRewriteHttpDesc,
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"shared_cookies": {
							Description: proxySharedCookies,
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
			"rdp": {
				Description: rdpDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remote_app": {
							Description: rdpRemoteAppDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"remote_app_cmd_args": {
							Description: rdpRemoteAppCmdArgsDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"remote_app_work_dir": {
							Description: rdpRemoteAppWorkDirDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"security": {
							Description: rdpRemoteAppSecurityDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
						"server_keyboard_layout": {
							Description: rdpSeverKeyboardLayoutDesc,
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
			"version": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}
