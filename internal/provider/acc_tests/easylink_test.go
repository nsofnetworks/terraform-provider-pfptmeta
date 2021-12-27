package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	easylinkDependencies = `
resource "pfptmeta_group" "new_group" {
  name = "easylink-group"
}

locals {
  hostname = "cert-test.notariustest.com"
  ipv4     = "196.10.10.1"
}

resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service name"
  mapped_service = local.ipv4
}

resource "pfptmeta_network_element_alias" "alias" {
  network_element_id = pfptmeta_network_element.mapped-service.id
  alias              = local.hostname
}
`
	certificateDependency = `
resource "pfptmeta_certificate" "cert" {
  name        = "certificate name"
  sans        = [local.hostname]
}
`
	metaEasylinkStep1 = `
resource "pfptmeta_easylink" "meta_easylink" {
  name        = "meta easylink name"
  description = "meta easylink description"
  domain_name = local.hostname
  access_type = "meta"
  port        = 443
  protocol    = "https"
  viewers     = [pfptmeta_group.new_group.id]
}`
	metaEasylinkStep2 = `
resource "pfptmeta_easylink" "meta_easylink" {
  name        = "meta easylink name1"
  description = "meta easylink description1"
  domain_name = local.hostname
  access_type = "meta"
  port        = 443
  protocol    = "https"
  viewers     = [pfptmeta_group.new_group.id]
}`
	metaRdpEasylinkStep1 = `
resource "pfptmeta_easylink" "meta_rdp_easylink" {
  name        = "meta_rdp easylink name"
  domain_name = local.hostname
  access_type = "meta"
  port        = 3389
  protocol    = "rdp"
  viewers     = [pfptmeta_group.new_group.id]
  rdp {
    security               = "nla"
    server_keyboard_layout = "english-us"
  }
}
`
	metaRdpEasylinkStep2 = `
resource "pfptmeta_easylink" "meta_rdp_easylink" {
  name        = "meta_rdp easylink name1"
  domain_name = local.hostname
  access_type = "meta"
  port        = 3389
  protocol    = "rdp"
  viewers     = [pfptmeta_group.new_group.id]
  rdp {
    security               = "rdp"
    server_keyboard_layout = "italian"
  }
}`
	redirectEasylinkStep = `
resource "pfptmeta_easylink" "redirect_easylink" {
  name              = "redirect easylink name"
  domain_name       = local.ipv4
  access_fqdn       = local.hostname
  access_type       = "redirect"
  port              = 443
  protocol          = "https"
  mapped_element_id = pfptmeta_network_element.mapped-service.id
  viewers           = [pfptmeta_group.new_group.id]
  certificate_id    = pfptmeta_certificate.cert.id
  root_path         = "/application"
}`
	nativeProxyEasylinkStep1 = `
resource "pfptmeta_easylink" "native_easylink" {
  name              = "native easylink name"
  domain_name       = local.ipv4
  access_fqdn       = local.hostname
  access_type       = "native"
  port              = 443
  protocol          = "https"
  mapped_element_id = pfptmeta_network_element.mapped-service.id
  viewers           = [pfptmeta_group.new_group.id]
  certificate_id    = pfptmeta_certificate.cert.id
  root_path         = "/application"
  proxy {
    rewrite_content_types = ["json"]
    rewrite_http          = true
    rewrite_hosts         = true
  }
}
`
	nativeProxyEasylinkStep2 = `
resource "pfptmeta_easylink" "native_easylink" {
  name              = "native easylink name1"
  domain_name       = local.ipv4
  access_fqdn       = local.hostname
  access_type       = "native"
  port              = 443
  protocol          = "https"
  mapped_element_id = pfptmeta_network_element.mapped-service.id
  viewers           = [pfptmeta_group.new_group.id]
  certificate_id    = pfptmeta_certificate.cert.id
  root_path         = "/application1"
  proxy {
    rewrite_content_types = ["javascript"]
    rewrite_http          = false
    rewrite_hosts         = true
  }
}
`
	dataSourceEasylink = `
data "pfptmeta_easylink" "easylink" {
  id = pfptmeta_easylink.meta_easylink.id
}
`
)

func TestAccMetaEasylink(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("easylink", "v1/easylinks"),
		Steps: []resource.TestStep{
			{
				Config: easylinkDependencies + metaEasylinkStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_easylink.meta_easylink", "id", regexp.MustCompile("^el-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_easylink", "name", "meta easylink name"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_easylink", "description", "meta easylink description"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_easylink", "domain_name", "cert-test.notariustest.com"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_easylink", "access_type", "meta"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_easylink", "port", "443"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_easylink", "protocol", "https"),
					resource.TestMatchResourceAttr("pfptmeta_easylink.meta_easylink", "viewers.0", regexp.MustCompile("^grp-.+$")),
				),
			},
			{
				Config: easylinkDependencies + metaEasylinkStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_easylink.meta_easylink", "id", regexp.MustCompile("^el-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_easylink", "name", "meta easylink name1"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_easylink", "description", "meta easylink description1"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_easylink", "domain_name", "cert-test.notariustest.com"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_easylink", "access_type", "meta"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_easylink", "port", "443"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_easylink", "protocol", "https"),
					resource.TestMatchResourceAttr("pfptmeta_easylink.meta_easylink", "viewers.0", regexp.MustCompile("^grp-.+$")),
				),
			},
		},
	})
}

func TestAccMetaRdpEasylink(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("easylink", "v1/easylinks"),
		Steps: []resource.TestStep{
			{
				Config: easylinkDependencies + metaRdpEasylinkStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "id", regexp.MustCompile("^el-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "name", "meta_rdp easylink name"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "domain_name", "cert-test.notariustest.com"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "access_type", "meta"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "port", "3389"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "protocol", "rdp"),
					resource.TestMatchResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "viewers.0", regexp.MustCompile("^grp-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "rdp.0.security", "nla"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "rdp.0.server_keyboard_layout", "english-us"),
				),
			},
			{
				Config: easylinkDependencies + metaRdpEasylinkStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "id", regexp.MustCompile("^el-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "name", "meta_rdp easylink name1"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "domain_name", "cert-test.notariustest.com"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "access_type", "meta"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "port", "3389"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "protocol", "rdp"),
					resource.TestMatchResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "viewers.0", regexp.MustCompile("^grp-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "rdp.0.security", "rdp"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.meta_rdp_easylink", "rdp.0.server_keyboard_layout", "italian"),
				),
			},
		},
	})
}

func TestAccNativeAndRedirectEasylink(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccReleasePreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("easylink", "v1/easylinks"),
		Steps: []resource.TestStep{
			{
				Config: easylinkDependencies + certificateDependency + redirectEasylinkStep,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_easylink.redirect_easylink", "id", regexp.MustCompile("^el-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_easylink.redirect_easylink", "name", "redirect easylink name"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.redirect_easylink", "domain_name", "196.10.10.1"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.redirect_easylink", "access_fqdn", "cert-test.notariustest.com"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.redirect_easylink", "access_type", "redirect"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.redirect_easylink", "port", "443"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.redirect_easylink", "protocol", "https"),
					resource.TestMatchResourceAttr("pfptmeta_easylink.redirect_easylink", "viewers.0", regexp.MustCompile("^grp-.+$")),
					resource.TestMatchResourceAttr("pfptmeta_easylink.redirect_easylink", "certificate_id", regexp.MustCompile("^crt-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_easylink.redirect_easylink", "root_path", "/application"),
				),
			},
			{
				Config: easylinkDependencies + certificateDependency + nativeProxyEasylinkStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_easylink.native_easylink", "id", regexp.MustCompile("^el-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "name", "native easylink name"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "domain_name", "196.10.10.1"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "access_fqdn", "cert-test.notariustest.com"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "access_type", "native"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "port", "443"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "protocol", "https"),
					resource.TestMatchResourceAttr("pfptmeta_easylink.native_easylink", "viewers.0", regexp.MustCompile("^grp-.+$")),
					resource.TestMatchResourceAttr("pfptmeta_easylink.native_easylink", "certificate_id", regexp.MustCompile("^crt-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "root_path", "/application"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "proxy.0.rewrite_content_types.0", "json"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "proxy.0.rewrite_http", "true"),
				),
			},
			{
				Config: easylinkDependencies + certificateDependency + nativeProxyEasylinkStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_easylink.native_easylink", "id", regexp.MustCompile("^el-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "name", "native easylink name1"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "domain_name", "196.10.10.1"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "access_fqdn", "cert-test.notariustest.com"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "access_type", "native"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "port", "443"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "protocol", "https"),
					resource.TestMatchResourceAttr("pfptmeta_easylink.native_easylink", "viewers.0", regexp.MustCompile("^grp-.+$")),
					resource.TestMatchResourceAttr("pfptmeta_easylink.native_easylink", "certificate_id", regexp.MustCompile("^crt-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "root_path", "/application1"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "proxy.0.rewrite_content_types.0", "javascript"),
					resource.TestCheckResourceAttr("pfptmeta_easylink.native_easylink", "proxy.0.rewrite_http", "false"),
				),
			},
		},
	})
}

func TestAccMDataSourceEasylink(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("easylink", "v1/easylinks"),
		Steps: []resource.TestStep{
			{
				Config: easylinkDependencies + metaEasylinkStep1 + dataSourceEasylink,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_easylink.easylink", "id", regexp.MustCompile("^el-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_easylink.easylink", "name", "meta easylink name"),
					resource.TestCheckResourceAttr("data.pfptmeta_easylink.easylink", "description", "meta easylink description"),
					resource.TestCheckResourceAttr("data.pfptmeta_easylink.easylink", "domain_name", "cert-test.notariustest.com"),
					resource.TestCheckResourceAttr("data.pfptmeta_easylink.easylink", "access_type", "meta"),
					resource.TestCheckResourceAttr("data.pfptmeta_easylink.easylink", "port", "443"),
					resource.TestCheckResourceAttr("data.pfptmeta_easylink.easylink", "protocol", "https"),
					resource.TestMatchResourceAttr("data.pfptmeta_easylink.easylink", "viewers.0", regexp.MustCompile("^grp-.+$")),
				),
			},
		},
	})
}
