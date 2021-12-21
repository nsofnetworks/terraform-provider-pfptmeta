package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	certificateResourceStep1 = `
resource "pfptmeta_certificate" "cert" {
  name = "c-name"
  description = "c-description"
  sans = ["cert-test.notariustest.com"]
}`
	certificateResourceStep2 = `
resource "pfptmeta_certificate" "cert" {
  name = "c-name1"
  description = "c-description1"
  sans = ["cert-test.notariustest.com"]
}`
	certificateDataSource = `

data "pfptmeta_certificate" "cert" {
  id = pfptmeta_certificate.cert.id
}`
)

func TestAccCertificate(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("certificate", "v1/certificates"),
		Steps: []resource.TestStep{
			{
				Config: certificateResourceStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_certificate.cert", "id", regexp.MustCompile("^crt-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_certificate.cert", "name", "c-name"),
					resource.TestCheckResourceAttr("pfptmeta_certificate.cert", "description", "c-description"),
					resource.TestCheckResourceAttr("pfptmeta_certificate.cert", "sans.0", "cert-test.notariustest.com"),
				),
			},
			{
				Config: certificateResourceStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_certificate.cert", "id", regexp.MustCompile("^crt-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_certificate.cert", "name", "c-name1"),
					resource.TestCheckResourceAttr("pfptmeta_certificate.cert", "description", "c-description1"),
					resource.TestCheckResourceAttr("pfptmeta_certificate.cert", "sans.0", "cert-test.notariustest.com"),
				),
			},
			{
				Config: certificateResourceStep1 + certificateDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.pfptmeta_certificate.cert", "id", regexp.MustCompile("^crt-.+$")),
					resource.TestCheckResourceAttr("data.pfptmeta_certificate.cert", "name", "c-name"),
					resource.TestCheckResourceAttr("data.pfptmeta_certificate.cert", "description", "c-description"),
					resource.TestCheckResourceAttr("data.pfptmeta_certificate.cert", "sans.0", "cert-test.notariustest.com"),
				),
			},
		},
	})
}
