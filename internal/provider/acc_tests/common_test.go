package acc_tests

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	p "github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider"
	"net/http"
	"os"
	"testing"
)

var provider *schema.Provider

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"pfptmeta": func() (*schema.Provider, error) {
		if provider == nil {
			provider = p.New("dev")()
		}
		return provider, nil
	},
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("PFPTMETA_API_KEY") == "" {
		t.Fatalf("PFPTMETA_API_KEY env var must be set")
	}
	if os.Getenv("PFPTMETA_API_SECRET") == "" {
		t.Fatalf("PFPTMETA_API_SECRET env var must be set")
	}
	if os.Getenv("PFPTMETA_ORG_SHORTNAME") == "" {
		t.Fatalf("PFPTMETA_ORG_SHORTNAME env var must be set")
	}

}

func validateResourceDestroyed(resource, resourcePath string) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		c := provider.Meta().(*client.Client)
		resourceType := fmt.Sprintf("pfptmeta_%s", resource)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != resourceType {
				continue
			}
			neId := rs.Primary.ID
			_, err := c.GetResource(resourcePath, neId)
			if err == nil {
				return fmt.Errorf("%s %s still exists", resource, neId)
			}
			errResponse, ok := err.(*client.ErrorResponse)
			if ok && errResponse.Status == http.StatusNotFound {
				return nil
			}
			return fmt.Errorf("failed to verify %s %s was destroyed: %s", resource, neId, err)
		}
		return nil
	}
}
