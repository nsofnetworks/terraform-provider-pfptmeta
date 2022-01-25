package log_streaming_access_bridge

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
	"regexp"
)

func Resource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		CreateContext: abCreate,
		ReadContext:   abRead,
		UpdateContext: abUpdate,
		DeleteContext: abDelete,
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
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"notification_channels": {
				Description: notificationChannelDesc,
				Type:        schema.TypeList,
				MaxItems:    3,
				Optional:    true,
				Elem: &schema.Schema{
					Type:             schema.TypeString,
					ValidateDiagFunc: common.ValidateID(false, "nch"),
				},
			},
			"siem_config": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"export_logs": {
							Description: exportLogsDesc,
							Type:        schema.TypeList,
							Required:    true,
							MinItems:    1,
							Elem: &schema.Schema{
								Type:             schema.TypeString,
								ValidateDiagFunc: common.ValidateStringENUM("api", "traffic", "security", "metaproxy", "webfilter"),
							},
						},
						"proofpoint_casb_config": {
							Description: proofpointCASBConfig,
							Type:        schema.TypeList,
							Optional:    true,
							ConflictsWith: []string{
								"siem_config.0.qradar_http_config",
								"siem_config.0.s3_config",
								"siem_config.0.splunk_http_config",
								"siem_config.0.syslog_config",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Description:      casbRegionDesc,
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: common.ValidateStringENUM("US", "EU"),
									},
									"tenant_id": {
										Description:      casbTenantIDDesc,
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: common.ValidatePattern(regexp.MustCompile("^tenant_([a-f0-9]+)$")),
									},
								},
							},
						},
						"qradar_http_config": {
							Description: qRadarConfDesc,
							Type:        schema.TypeList,
							Optional:    true,
							ConflictsWith: []string{
								"siem_config.0.proofpoint_casb_config",
								"siem_config.0.s3_config",
								"siem_config.0.splunk_http_config",
								"siem_config.0.syslog_config",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate": {
										Description:      qRadarCertDesc,
										Type:             schema.TypeString,
										Optional:         true,
										Sensitive:        true,
										ValidateDiagFunc: common.ValidatePEMCert(),
									},
									"url": {
										Description:      qRadarURLDesc,
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: common.ComposeOrValidations(common.ValidateURL(), common.ValidateHostName()),
									},
								},
							},
						},
						"s3_config": {
							Description: s3ConfDesc,
							Type:        schema.TypeList,
							Optional:    true,
							ConflictsWith: []string{
								"siem_config.0.proofpoint_casb_config",
								"siem_config.0.qradar_http_config",
								"siem_config.0.splunk_http_config",
								"siem_config.0.syslog_config",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Description: s3BucketDesc,
										Type:        schema.TypeString,
										Required:    true,
									},
									"compress": {
										Description: s3CompressDesc,
										Type:        schema.TypeBool,
										Required:    true,
									},
									"prefix": {
										Description: s3Prefix,
										Type:        schema.TypeString,
										Optional:    true,
									},
								},
							},
						},
						"splunk_http_config": {
							Description: splunkDesc,
							Type:        schema.TypeList,
							Optional:    true,
							ConflictsWith: []string{
								"siem_config.0.proofpoint_casb_config",
								"siem_config.0.qradar_http_config",
								"siem_config.0.s3_config",
								"siem_config.0.syslog_config",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate": {
										Description:      splunkCerDesc,
										Type:             schema.TypeString,
										Optional:         true,
										Sensitive:        true,
										ValidateDiagFunc: common.ValidatePEMCert(),
									},
									"publicly_accessible": {
										Description: splunkPubliclyAccessibleDesc,
										Type:        schema.TypeBool,
										Optional:    true,
									},
									"token": {
										Description: splunkTokenDesc,
										Type:        schema.TypeString,
										Required:    true,
										Sensitive:   true,
									},
									"url": {
										Description:      splunkURLDesc,
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: common.ComposeOrValidations(common.ValidateURL(), common.ValidateHostName()),
									},
								},
							},
						},
						"syslog_config": {
							Description: syslogDesc,
							Type:        schema.TypeList,
							Optional:    true,
							ConflictsWith: []string{
								"siem_config.0.proofpoint_casb_config",
								"siem_config.0.qradar_http_config",
								"siem_config.0.s3_config",
								"siem_config.0.splunk_http_config",
							},
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Description:      syslogHostDesc,
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: common.ValidateHostName(),
									},
									"port": {
										Description:      syslogPortDesc,
										Type:             schema.TypeInt,
										Required:         true,
										ValidateDiagFunc: common.ValidateIntRange(1, 65535),
									},
									"proto": {
										Description:      syslogProtoDesc,
										Type:             schema.TypeString,
										Required:         true,
										ValidateDiagFunc: common.ValidateStringENUM("tcp", "udp"),
									},
								},
							},
						},
					},
				},
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
