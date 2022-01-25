package log_streaming_access_bridge

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
)

func DataSource() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: description,

		ReadContext: abRead,

		Schema: map[string]*schema.Schema{
			"id": {
				Type:             schema.TypeString,
				Required:         true,
				ValidateDiagFunc: common.ValidateID(false, "ab"),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"notification_channels": {
				Description: notificationChannelDesc,
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"siem_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"export_logs": {
							Description: exportLogsDesc,
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
						},
						"proofpoint_casb_config": {
							Description: proofpointCASBConfig,
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"region": {
										Description: casbRegionDesc,
										Type:        schema.TypeString,
										Computed:    true,
									},
									"tenant_id": {
										Description: casbTenantIDDesc,
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"qradar_http_config": {
							Description: qRadarConfDesc,
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate": {
										Description: qRadarCertDesc,
										Type:        schema.TypeString,
										Computed:    true,
										Sensitive:   true,
									},
									"url": {
										Description: qRadarURLDesc,
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"s3_config": {
							Description: s3ConfDesc,
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"bucket": {
										Description: s3BucketDesc,
										Type:        schema.TypeString,
										Computed:    true,
									},
									"compress": {
										Description: s3CompressDesc,
										Type:        schema.TypeBool,
										Computed:    true,
									},
									"prefix": {
										Description: s3Prefix,
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"splunk_http_config": {
							Description: splunkDesc,
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"certificate": {
										Description: splunkCerDesc,
										Type:        schema.TypeString,
										Computed:    true,
										Sensitive:   true,
									},
									"publicly_accessible": {
										Description: splunkPubliclyAccessibleDesc,
										Type:        schema.TypeBool,
										Computed:    true,
									},
									"token": {
										Description: splunkTokenDesc,
										Type:        schema.TypeString,
										Computed:    true,
										Sensitive:   true,
									},
									"url": {
										Description: splunkURLDesc,
										Type:        schema.TypeString,
										Computed:    true,
									},
								},
							},
						},
						"syslog_config": {
							Description: syslogDesc,
							Type:        schema.TypeList,
							Computed:    true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"host": {
										Description: syslogHostDesc,
										Type:        schema.TypeString,
										Computed:    true,
									},
									"port": {
										Description: syslogPortDesc,
										Type:        schema.TypeInt,
										Computed:    true,
									},
									"proto": {
										Description: syslogProtoDesc,
										Type:        schema.TypeString,
										Computed:    true,
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
