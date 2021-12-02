package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceEnterpriseDNS(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		CheckDestroy:      validateResourceDestroyed("enterprise_dns", "v1/enterprise_dns"),
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceEnterpriseDNSStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_enterprise_dns.enterprise_dns", "id", regexp.MustCompile("^ed-[\\d]+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_enterprise_dns.enterprise_dns", "name", "ed-name",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_enterprise_dns.enterprise_dns", "description", "ed-description",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_enterprise_dns.enterprise_dns", "mapped_domains.0.name", "step1.test1.com",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_enterprise_dns.enterprise_dns", "mapped_domains.0.mapped_domain", "step1.test1.com",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_enterprise_dns.enterprise_dns", "mapped_domains.1.name", "step1.test2.com",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_enterprise_dns.enterprise_dns", "mapped_domains.1.mapped_domain", "step1.test2.com",
					),
				),
			},
			{
				Config: testAccResourceEnterpriseDNSStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pfptmeta_enterprise_dns.enterprise_dns", "name", "ed-name1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_enterprise_dns.enterprise_dns", "description", "ed-description1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_enterprise_dns.enterprise_dns", "mapped_domains.0.name", "step2.test1.com",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_enterprise_dns.enterprise_dns", "mapped_domains.0.mapped_domain", "step2.test1.com",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_enterprise_dns.enterprise_dns", "mapped_domains.1.name", "step2.test2.com",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_enterprise_dns.enterprise_dns", "mapped_domains.1.mapped_domain", "step2.test2.com",
					),
				),
			},
		},
	})
}

const testAccResourceEnterpriseDNSStep1 = `
resource "pfptmeta_enterprise_dns" "enterprise_dns" {
  name           = "ed-name"
  description    = "ed-description"
  mapped_domains {
    name          = "step1.test1.com"
    mapped_domain = "step1.test1.com"
  }
  mapped_domains {
    name          = "step1.test2.com"
    mapped_domain = "step1.test2.com"
  }
}
`

const testAccResourceEnterpriseDNSStep2 = `
resource "pfptmeta_enterprise_dns" "enterprise_dns" {
  name           = "ed-name1"
  description    = "ed-description1"
  mapped_domains {
    name          = "step2.test1.com"
    mapped_domain = "step2.test1.com"
  }
  mapped_domains {
    name          = "step2.test2.com"
    mapped_domain = "step2.test2.com"
  }
}
`
