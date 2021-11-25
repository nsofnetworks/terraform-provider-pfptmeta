package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
	"regexp"
	"testing"
)

func TestAccResourceProtocolGroup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      checkProtocolGroupsDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceProtocolGroup,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(
						"pfptmeta_protocol_group.new_protocol", "id", regexp.MustCompile("^pg-.*$"),
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_protocol_group.new_protocol", "name", "NEW_PROTOCOL",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_protocol_group.new_protocol", "read_only", "false",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_protocol_group.new_protocol", "protocols.0.from_port", "445",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_protocol_group.new_protocol", "protocols.0.to_port", "445",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_protocol_group.new_protocol", "protocols.0.proto", "udp",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_protocol_group.new_protocol", "protocols.1.from_port", "446",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_protocol_group.new_protocol", "protocols.1.to_port", "446",
					),
					resource.TestCheckResourceAttr(
						"pfptmeta_protocol_group.new_protocol", "protocols.1.proto", "tcp",
					),
				),
			},
		},
	})
}

func checkProtocolGroupsDestroyed(s *terraform.State) error {
	c := provider.Meta().(*client.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "pfptmeta_protocol_group" {
			continue
		}
		pID := rs.Primary.ID
		_, err := client.GetProtocolGroupById(c, pID)
		if err == nil {
			return fmt.Errorf("protocol group %s still exists", pID)
		}
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			return nil
		}
		return fmt.Errorf("failed to verify protocol group %s was destroyed: %s", pID, err)
	}

	return nil
}

const testAccResourceProtocolGroup = `
resource "pfptmeta_protocol_group" "new_protocol" {
  name = "NEW_PROTOCOL"
  protocols {
    from_port = 445
    to_port   = 445
    proto     = "udp"
  }
  protocols {
    from_port = 446
    to_port   = 446
    proto     = "tcp"
  }
}
`
