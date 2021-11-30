package provider

import (
	"regexp"
	"testing"
)

func TestValidateENUM(t *testing.T) {
	cases := map[string]struct {
		Input       string
		Enum        []string
		ShouldError bool
	}{
		"positive-test": {
			Input:       "test1",
			Enum:        []string{"test3", "test2", "test1"},
			ShouldError: false,
		},
		"negative-test": {
			Input:       "test4",
			Enum:        []string{"test3", "test2", "test1"},
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

func TestValidateID(t *testing.T) {
	cases := map[string]struct {
		Input       string
		numeric     bool
		prefixes    []string
		ShouldError bool
	}{
		"positive-test1": {
			Input:       "ne-123",
			numeric:     true,
			prefixes:    []string{"ne"},
			ShouldError: false,
		},
		"positive-test2": {
			Input:       "ne-123",
			numeric:     true,
			prefixes:    []string{"ed", "ne"},
			ShouldError: false,
		},
		"negative-test-no-id": {
			Input:       "12345",
			prefixes:    []string{"ne"},
			ShouldError: true,
		},
		"negative-test-id-with-no-suffix": {
			Input:       "ne-",
			prefixes:    []string{"ne"},
			ShouldError: true,
		},
		"negative-test-non-numeric": {
			Input:       "ne-abc",
			numeric:     true,
			prefixes:    []string{"ne"},
			ShouldError: true,
		},
		"negative-test-non-alphabetic": {
			Input:       "ne-!@#",
			numeric:     true,
			prefixes:    []string{"ne"},
			ShouldError: true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			diags := validateID(tc.numeric, tc.prefixes...)(tc.Input, nil)
			if diags.HasError() && !tc.ShouldError {
				t.Errorf("%s failed: %+v", name, diags[0])
			}
		})
	}
}

func TestValidatePattern(t *testing.T) {
	cases := map[string]struct {
		Input       string
		pattern     *regexp.Regexp
		ShouldError bool
	}{
		"positive-test": {
			Input:       "test123",
			pattern:    regexp.MustCompile("test[\\d]+"),
			ShouldError: false,
		},
		"negative-test": {
			Input:       "abcd",
			pattern:    regexp.MustCompile("[\\d]+"),
			ShouldError: true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			diags := validatePattern(tc.pattern)(tc.Input, nil)
			if diags.HasError() && !tc.ShouldError {
				t.Errorf("%s failed: %+v", name, diags[0])
			}
		})
	}
}