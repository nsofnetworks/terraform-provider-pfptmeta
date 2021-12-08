package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/enterprise_dns"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/group"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/mapped_domain"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/mapped_host"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/metaport"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/metaport_cluster"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/metaport_failover"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/network_element"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/network_element_alias"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/protocol_group"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/role"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/routing_group"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/user"
)

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		p := &schema.Provider{
			Schema: map[string]*schema.Schema{
				"api_key": {
					Description: "Alternatively use the `PFPTMETA_API_KEY` env variable",
					Type:        schema.TypeString,
					DefaultFunc: schema.EnvDefaultFunc("PFPTMETA_API_KEY", nil),
					Optional:    true,
				},
				"api_secret": {
					Description: "Alternatively use the `PFPTMETA_API_SECRET` env variable",
					Type:        schema.TypeString,
					DefaultFunc: schema.EnvDefaultFunc("PFPTMETA_API_SECRET", nil),
					Optional:    true,
					Sensitive:   true,
				},
				"org_shortname": {
					Description: "Alternatively use the `PFPTMETA_ORG_SHORTNAME` env variable",
					Type:        schema.TypeString,
					DefaultFunc: schema.EnvDefaultFunc("PFPTMETA_ORG_SHORTNAME", nil),
					Optional:    true,
				},
			},
			DataSourcesMap: map[string]*schema.Resource{
				"pfptmeta_network_element_alias": network_element_alias.DataSource(),
				"pfptmeta_mapped_domain":         mapped_domain.DataSource(),
				"pfptmeta_mapped_host":           mapped_host.DataSource(),
				"pfptmeta_network_element":       network_element.DataSource(),
				"pfptmeta_metaport":              metaport.DataSource(),
				"pfptmeta_metaport_cluster":      metaport_cluster.DataSource(),
				"pfptmeta_metaport_failover":     metaport_failover.DataSource(),
				"pfptmeta_enterprise_dns":        enterprise_dns.DataSource(),
				"pfptmeta_protocol_group":        protocol_group.DataSource(),
				"pfptmeta_role":                  role.DataSource(),
				"pfptmeta_group":                 group.DataSource(),
				"pfptmeta_user":                  user.DataSource(),
				"pfptmeta_routing_group":         routing_group.DataSource(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"pfptmeta_network_element":       network_element.Resource(),
				"pfptmeta_network_element_alias": network_element_alias.Resource(),
				"pfptmeta_mapped_domain":         mapped_domain.Resource(),
				"pfptmeta_mapped_host":           mapped_host.Resource(),
				"pfptmeta_metaport":              metaport.Resource(),
				"pfptmeta_metaport_cluster":      metaport_cluster.Resource(),
				"pfptmeta_metaport_failover":     metaport_failover.Resource(),
				"pfptmeta_enterprise_dns":        enterprise_dns.Resource(),
				"pfptmeta_protocol_group":        protocol_group.Resource(),
				"pfptmeta_role":                  role.Resource(),
				"pfptmeta_group":                 group.Resource(),
				"pfptmeta_user":                  user.Resource(),
				"pfptmeta_routing_group":         routing_group.Resource(),
			},
		}
		p.ConfigureContextFunc = configure(version, p)
		return p
	}
}

func configure(version string, p *schema.Provider) func(context.Context, *schema.ResourceData) (interface{}, diag.Diagnostics) {
	return func(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		userAgent := p.UserAgent("terraform-provider-pfptmeta", version)
		c, err := client.NewClient(d, userAgent)
		if err != nil {
			return nil, diag.FromErr(err)
		}
		return c, nil
	}
}
