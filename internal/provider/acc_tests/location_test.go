package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccDataSourceLocation(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: AccTestLocation,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.pfptmeta_location.new_york", "city", expectedCity),
					resource.TestCheckResourceAttr("data.pfptmeta_location.new_york", "country", expectedCountry),
					resource.TestCheckResourceAttr("data.pfptmeta_location.new_york", "name", expectedName),
					resource.TestCheckResourceAttr("data.pfptmeta_location.new_york", "state", expectedState),
				),
			},
		},
	})
}

const (
	expectedCity    = "New York"
	expectedCountry = "United States"
	expectedName    = "LGA"
	expectedState   = "New York"
)

const AccTestLocation = `
data "pfptmeta_location" "new_york" {
  name = "LGA"
}
`
