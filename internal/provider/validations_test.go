package provider

import (
	"testing"
)

func TestValidateENUM(t *testing.T) {
	cases := map[string]struct {
		Input       string
		Enum        []string
		ShouldError bool
	}{
		"positive-test": {
			Input: "test1",
			Enum:  []string{"test3", "test2", "test1"},
			ShouldError: false,
		},
		"negative-test": {
			Input: "test4",
			Enum:  []string{"test3", "test2", "test1"},
			ShouldError: true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			diags := validateENUM(tc.Enum...)(tc.Input, nil)
			if diags.HasError() && !tc.ShouldError {
				t.Errorf("%s failed: %+v", name, diags[0])
			}
		})
	}
}
