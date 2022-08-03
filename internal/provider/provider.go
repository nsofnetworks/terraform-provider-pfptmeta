package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/access_control"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/alert"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/certificate"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/device"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/device_alias"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/device_settings"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/easylink"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/egress_route"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/enterprise_dns"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/group"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/group_roles_attachment"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/group_users_attachment"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/location"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/log_streaming_access_bridge"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/mapped_domain"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/mapped_host"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/metaport"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/metaport_cluster"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/metaport_failover"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/metaport_mapped_elements_attachment"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/network_element"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/network_element_alias"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/notification_channel"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/policy"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/posture_check"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/protocol_group"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/role"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/routing_group"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/routing_group_mapped_elements_attachment"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/trusted_network"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/user"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/user_roles_attachment"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/user_settings"
)

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_key": {
					Description:      "Alternatively, use the `PFPTMETA_API_KEY` env variable",
					Type:             schema.TypeString,
					DefaultFunc:      schema.EnvDefaultFunc("PFPTMETA_API_KEY", nil),
					Optional:         true,
					ValidateDiagFunc: common.ValidateID(false, "key"),
				},
				"api_secret": {
					Description: "Alternatively, use the `PFPTMETA_API_SECRET` env variable",
					Type:        schema.TypeString,
					DefaultFunc: schema.EnvDefaultFunc("PFPTMETA_API_SECRET", nil),
					Optional:    true,
					Sensitive:   true,
				},
				"org_shortname": {
					Description: "Alternatively, use the `PFPTMETA_ORG_SHORTNAME` env variable",
					Type:        schema.TypeString,
					DefaultFunc: schema.EnvDefaultFunc("PFPTMETA_ORG_SHORTNAME", nil),
					Optional:    true,
				},
				"realm": {
					Description:      "GDPR data location, ENUM: `us`, `eu`. defaults to `us`",
					Type:             schema.TypeString,
					DefaultFunc:      schema.EnvDefaultFunc("PFPTMETA_REALM", nil),
					Optional:         true,
					ValidateDiagFunc: common.ValidateStringENUM("eu", "us"),
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"pfptmeta_network_element_alias":       network_element_alias.DataSource(),
				"pfptmeta_device_alias":                device_alias.DataSource(),
				"pfptmeta_mapped_domain":               mapped_domain.DataSource(),
				"pfptmeta_mapped_host":                 mapped_host.DataSource(),
				"pfptmeta_network_element":             network_element.DataSource(),
				"pfptmeta_device":                      device.DataSource(),
				"pfptmeta_metaport":                    metaport.DataSource(),
				"pfptmeta_metaport_cluster":            metaport_cluster.DataSource(),
				"pfptmeta_metaport_failover":           metaport_failover.DataSource(),
				"pfptmeta_enterprise_dns":              enterprise_dns.DataSource(),
				"pfptmeta_protocol_group":              protocol_group.DataSource(),
				"pfptmeta_role":                        role.DataSource(),
				"pfptmeta_group":                       group.DataSource(),
				"pfptmeta_user":                        user.DataSource(),
				"pfptmeta_notification_channel":        notification_channel.DataSource(),
				"pfptmeta_routing_group":               routing_group.DataSource(),
				"pfptmeta_policy":                      policy.DataSource(),
				"pfptmeta_location":                    location.DataSource(),
				"pfptmeta_egress_route":                egress_route.DataSource(),
				"pfptmeta_alert":                       alert.DataSource(),
				"pfptmeta_certificate":                 certificate.DataSource(),
				"pfptmeta_easylink":                    easylink.DataSource(),
				"pfptmeta_posture_check":               posture_check.DataSource(),
				"pfptmeta_access_control":              access_control.DataSource(),
				"pfptmeta_device_settings":             device_settings.DataSource(),
				"pfptmeta_trusted_network":             trusted_network.DataSource(),
				"pfptmeta_user_settings":               user_settings.DataSource(),
				"pfptmeta_log_streaming_access_bridge": log_streaming_access_bridge.DataSource(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"pfptmeta_network_element":                          network_element.Resource(),
				"pfptmeta_device":                                   device.Resource(),
				"pfptmeta_network_element_alias":                    network_element_alias.Resource(),
				"pfptmeta_device_alias":                             device_alias.Resource(),
				"pfptmeta_mapped_domain":                            mapped_domain.Resource(),
				"pfptmeta_mapped_host":                              mapped_host.Resource(),
				"pfptmeta_metaport":                                 metaport.Resource(),
				"pfptmeta_metaport_mapped_elements_attachment":      metaport_mapped_elements_attachment.Resource(),
				"pfptmeta_metaport_cluster":                         metaport_cluster.Resource(),
				"pfptmeta_metaport_failover":                        metaport_failover.Resource(),
				"pfptmeta_enterprise_dns":                           enterprise_dns.Resource(),
				"pfptmeta_protocol_group":                           protocol_group.Resource(),
				"pfptmeta_role":                                     role.Resource(),
				"pfptmeta_group":                                    group.Resource(),
				"pfptmeta_user":                                     user.Resource(),
				"pfptmeta_routing_group":                            routing_group.Resource(),
				"pfptmeta_routing_group_mapped_elements_attachment": routing_group_mapped_elements_attachment.Resource(),
				"pfptmeta_policy":                                   policy.Resource(),
				"pfptmeta_group_roles_attachment":                   group_roles_attachment.Resource(),
				"pfptmeta_group_users_attachment":                   group_users_attachment.Resource(),
				"pfptmeta_notification_channel":                     notification_channel.Resource(),
				"pfptmeta_egress_route":                             egress_route.Resource(),
				"pfptmeta_alert":                                    alert.Resource(),
				"pfptmeta_certificate":                              certificate.Resource(),
				"pfptmeta_easylink":                                 easylink.Resource(),
				"pfptmeta_posture_check":                            posture_check.Resource(),
				"pfptmeta_access_control":                           access_control.Resource(),
				"pfptmeta_device_settings":                          device_settings.Resource(),
				"pfptmeta_trusted_network":                          trusted_network.Resource(),
				"pfptmeta_user_roles_attachment":                    user_roles_attachment.Resource(),
				"pfptmeta_user_settings":                            user_settings.Resource(),
				"pfptmeta_log_streaming_access_bridge":              log_streaming_access_bridge.Resource(),
			},
		}
		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		userAgent := p.UserAgent("terraform-provider-pfptmeta", version)
		c, err := client.NewClient(ctx, d, userAgent)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return c, nil
	}
}
