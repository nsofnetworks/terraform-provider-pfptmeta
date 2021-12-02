package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
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
				"pfptmeta_network_element_alias": dataSourceNetworkElementAlias(),
				"pfptmeta_mapped_domain":         dataSourceMappedDomain(),
				"pfptmeta_mapped_host":           dataSourceMappedHost(),
				"pfptmeta_network_element":       dataSourceNetworkElement(),
				"pfptmeta_metaport":              dataSourceMetaport(),
				"pfptmeta_metaport_cluster":      dataSourceMetaportCluster(),
				"pfptmeta_metaport_failover":     dataSourceMetaportFailover(),
				"pfptmeta_enterprise_dns":        dataSourceEnterpriseDNS(),
			},
			ResourcesMap: map[string]*schema.Resource{
				"pfptmeta_network_element":       resourceNetworkElement(),
				"pfptmeta_network_element_alias": resourceNetworkElementAlias(),
				"pfptmeta_mapped_domain":         resourceMappedDomain(),
				"pfptmeta_mapped_host":           resourceMappedHost(),
				"pfptmeta_metaport":              resourceMetaport(),
				"pfptmeta_metaport_cluster":      resourceMetaportCluster(),
				"pfptmeta_metaport_failover":     resourceMetaportFailover(),
				"pfptmeta_enterprise_dns":        resourceEnterpriseDNS(),
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
