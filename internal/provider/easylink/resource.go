package easylink

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
	"regexp"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description:   description,
		CreateContext: easyLinkCreate,
		ReadContext:   easyLinkRead,
		UpdateContext: easyLinkUpdate,
		DeleteContext: easyLinkDelete,
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
			"domain_name": {
				Description:      domainNameDesc,
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateHostnameOrIPV4(),
			},
			"viewers": {
				Description: viewersDesc,
				Type:        schema.TypeSet,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "usr", "grp"),
				},
				Required: true,
				MinItems: 1,
			},
			"access_fqdn": {
				Description:      accessFQDNDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateDNS(),
			},
			"access_type": {
				Description:      accessTypeDesc,
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateStringENUM("meta", "redirect", "native"),
			},
			"mapped_element_id": {
				Description:      mappedElementIDDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateID(true, "ne"),
			},
			"certificate_id": {
				Description:      certificateIDDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidateID(false, "crt"),
			},
			"audit": {
				Description: auditDesc,
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"enable_sni": {
				Description: enableSNIDesc,
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"port": {
				Description:      portDesc,
				Type:             schema.TypeInt,
				ValidateDiagFunc: common.ValidateIntRange(0, 65535),
				Required:         true,
			},
			"protocol": {
				Description:      protocolDesc,
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateStringENUM("ssh", "rdp", "vnc", "http", "https"),
				ForceNew:         true,
			},
			"root_path": {
				Description:      rootPathDesc,
				Type:             schema.TypeString,
				Optional:         true,
				ValidateDiagFunc: common.ValidatePattern(regexp.MustCompile("^/.*$")),
			},

			"proxy": {
				Description: proxyDesc,
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enterprise_access": {
							Description: proxyEnterpriseAccess,
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
						"hosts": {
							Description: proxyHostsDesc,
							Optional:    true,
							Computed:    true,
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidateHTTPNetLocation(),
							},
						},
						"http_host_header": {
							Description: proxyHostHeaderDesc,
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
						},
						"rewrite_content_types": {
							Description: proxyRewriteContentTypesDesc,
							Optional:    true,
							Computed:    true,
							Type:        schema.TypeList,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidateStringENUM("html", "json", "javascript", "text"),
							},
						},
						"rewrite_hosts": {
							Description: proxyRewriteHosts,
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
						"rewrite_hosts_client": {
							Description: proxyRewriteHostsClient,
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
						"rewrite_http": {
							Description: proxyRewriteHttpDesc,
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
						"shared_cookies": {
							Description: proxySharedCookies,
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"rdp": {
				Description: rdpDesc,
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"remote_app": {
							Description: rdpRemoteAppDesc,
							Type:        schema.TypeString,
							Optional:    true,
						},
						"remote_app_cmd_args": {
							Description: rdpRemoteAppCmdArgsDesc,
							Type:        schema.TypeString,
							Optional:    true,
						},
						"remote_app_work_dir": {
							Description: rdpRemoteAppWorkDirDesc,
							Type:        schema.TypeString,
							Optional:    true,
						},
						"security": {
							Description: rdpRemoteAppSecurityDesc,
							Type:        schema.TypeString,
							Optional:    true,
						},
						"server_keyboard_layout": {
							Description: rdpSeverKeyboardLayoutDesc,
							Type:        schema.TypeString,
							Optional:    true,
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
