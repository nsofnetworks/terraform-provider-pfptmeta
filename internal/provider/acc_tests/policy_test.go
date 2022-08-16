package acc_tests

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourcePolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("policy", "v1/policies"),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_policy.policy", "id", regexp.MustCompile("^pol-.+$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_policy.policy", "name", "policy name",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_policy.policy", "description", "policy description",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"pfptmeta_policy.policy", "sources.*", "pfptmeta_group.group", "id",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"pfptmeta_policy.policy", "destinations.*", "pfptmeta_network_element.mapped-service", "id",
					),
					resource.TestCheckTypeSetElemAttrPair(
						"pfptmeta_policy.policy", "protocol_groups.0", "data.pfptmeta_protocol_group.HTTPS", "id",
					),
				),
			},
			{
				Config: testAccPolicyStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"pfptmeta_policy.policy", "name", "policy name1",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_policy.policy", "description", "policy description1",
					),
					resource.TestCheckResourceAttr("pfptmeta_policy.policy", "sources.#", "0"),
					resource.TestCheckResourceAttr("pfptmeta_policy.policy", "destinations.#", "0"),
					resource.TestCheckResourceAttr("pfptmeta_policy.policy", "protocol_groups.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourcePolicy(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("policy", "v1/policies"),
		Steps: []resource.TestStep{
			{
				Config: testAccPolicyStep1 + policyDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"data.pfptmeta_policy.policy", "id", regexp.MustCompile("^pol-.+$"),
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_policy.policy", "name", "policy name",
					),
					resource.TestCheckResourceAttr(
						"data.pfptmeta_policy.policy", "description", "policy description",
					),
					resource.TestMatchResourceAttr(
						"data.pfptmeta_policy.policy", "sources.0", regexp.MustCompile("^grp-.+$"),
					),
					resource.TestMatchResourceAttr(
						"data.pfptmeta_policy.policy", "destinations.0", regexp.MustCompile("^ne-[\\d]+$"),
					),
					resource.TestMatchResourceAttr(
						"data.pfptmeta_policy.policy", "protocol_groups.0", regexp.MustCompile("^pg-.+$"),
					),
				),
			},
		},
	})
}

const testAccPolicyStep1 = `
resource "pfptmeta_group" "group" {
  name = "policy-test-group"
}

data "pfptmeta_protocol_group" "HTTPS" {
  name = "HTTPS"
}

resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service"
  mapped_service = "mapped.service.com"
}

resource "pfptmeta_policy" "policy" {
  name = "policy name"
  description = "policy description"
  sources = [pfptmeta_group.group.id]
  destinations = [pfptmeta_network_element.mapped-service.id]
  protocol_groups = [data.pfptmeta_protocol_group.HTTPS.id]
}
`

const testAccPolicyStep2 = `
resource "pfptmeta_group" "group" {
  name = "policy-test-group"
}

data "pfptmeta_protocol_group" "HTTPS" {
  name = "HTTPS"
}

resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service"
  mapped_service = "mapped.service.com"
}

resource "pfptmeta_policy" "policy" {
  name = "policy name1"
  description = "policy description1"
}
`

const policyDataSource = `

data "pfptmeta_policy" "policy" {
  id = pfptmeta_policy.policy.id
}`
