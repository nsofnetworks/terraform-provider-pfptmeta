package provider

import (
	"fmt"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateNetworkElementBody(t *testing.T) {
	cases := map[string]struct {
		Body          *client.NetworkElementBody
		Update        bool
		ExpectedError error
	}{
		"test_validation_with_owner_id_and_mapped_service": {
			Body: &client.NetworkElementBody{
				OwnerID:       "owner_id",
				MappedService: "service",
			},
			ExpectedError: fmt.Errorf("network element can only have one of \"mapped_subnets\", \"mapped_service\" or \"owner_id\""),
		},
		"test_validation_with_owner_id_and_mapped_subnet": {
			Body: &client.NetworkElementBody{
				OwnerID:       "owner_id",
				MappedSubnets: []string{"subnet"},
			},
			ExpectedError: fmt.Errorf("network element can only have one of \"mapped_subnets\", \"mapped_service\" or \"owner_id\""),
		},
		"test_validation_with_mapped_service_and_mapped_subnet": {
			Body: &client.NetworkElementBody{
				MappedSubnets: []string{"subnet"},
				MappedService: "service",
			},
			ExpectedError: fmt.Errorf("network element can only have one of \"mapped_subnets\", \"mapped_service\" or \"owner_id\""),
		},
		"test_validation_with_owner_id_on_update": {
			Body: &client.NetworkElementBody{
				OwnerID: "owner_id",
			},
			Update:        true,
			ExpectedError: fmt.Errorf("\"owner_id\" cannot be updated"),
		},
		"test_validation_with_platform_on_update": {
			Body: &client.NetworkElementBody{
				Platform: "platform",
			},
			Update:        true,
			ExpectedError: fmt.Errorf("\"platform\" cannot be updated"),
		},
		"test_validation_with_valid_body": {
			Body: &client.NetworkElementBody{},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			err := validateNetworkElementBody(tc.Body, tc.Update)
			assert.Equal(t, tc.ExpectedError, err)
		})
	}
}
