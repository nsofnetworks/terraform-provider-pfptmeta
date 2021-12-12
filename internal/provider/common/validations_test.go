package common

import (
	"github.com/stretchr/testify/assert"
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
			diags := ValidateENUM(tc.Enum...)(tc.Input, nil)
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
			diags := ValidateID(tc.numeric, tc.prefixes...)(tc.Input, nil)
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
			pattern:     regexp.MustCompile("test[\\d]+"),
			ShouldError: false,
		},
		"negative-test": {
			Input:       "abcd",
			pattern:     regexp.MustCompile("[\\d]+"),
			ShouldError: true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			diags := ValidatePattern(tc.pattern)(tc.Input, nil)
			if diags.HasError() && !tc.ShouldError {
				t.Errorf("%s failed: %+v", name, diags[0])
			}
		})
	}
}

func TestValidateHostname(t *testing.T) {
	cases := map[string]struct {
		Input       string
		ShouldError bool
	}{
		"positive-test": {
			Input:       "test.com",
			ShouldError: false,
		},
		"negative-test-numeric-tld": {
			Input:       "test.com.1234",
			ShouldError: true,
		},
		"negative-test-dot-suffix": {
			Input:       "test.com.",
			ShouldError: true,
		},
		"negative-test-hyphen-suffix": {
			Input:       "test-.com",
			ShouldError: true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			diags := ValidateHostName()(tc.Input, nil)
			if diags.HasError() && !tc.ShouldError {
				t.Errorf("%s failed: %+v", name, diags[0])
			}
		})
	}
}

func TestValidateIntRange(t *testing.T) {
	cases := map[string]struct {
		Input       int
		Min         int
		Max         int
		ShouldError bool
	}{
		"positive-test": {
			Input:       2,
			Min:         1,
			Max:         2,
			ShouldError: false,
		},
		"negative-test": {
			Input:       3,
			Min:         1,
			Max:         2,
			ShouldError: true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			diags := ValidateIntRange(tc.Min, tc.Max)(tc.Input, nil)
			if diags.HasError() && !tc.ShouldError {
				t.Errorf("%s failed: %+v", name, diags[0])
			}
		})
	}
}

func TestValidateHostnameOrIPV4(t *testing.T) {
	cases := map[string]struct {
		Input       string
		ShouldError bool
	}{
		"positive-test-hostname": {
			Input:       "test.com",
			ShouldError: false,
		},
		"positive-test-ipv4": {
			Input:       "127.0.0.1",
			ShouldError: false,
		},
		"negative-test-bad-ipv4": {
			Input:       "127.0.0.1.",
			ShouldError: true,
		},
		"negative-test-ipv6": {
			Input:       "2001:0db8:85a3:0000:0000:8a2e:0370:7334",
			ShouldError: true,
		},
		"negative-test-numeric-tld": {
			Input:       "test.com.1234",
			ShouldError: true,
		},
		"negative-test-dot-suffix": {
			Input:       "test.com.",
			ShouldError: true,
		},
		"negative-test-hyphen-suffix": {
			Input:       "test-.com",
			ShouldError: true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			diags := ValidateHostnameOrIPV4()(tc.Input, nil)
			if diags.HasError() && !tc.ShouldError {
				t.Errorf("%s failed: %+v", name, diags[0])
			}
		})
	}
}

var validPrivs = []string{
	"orgs:read",
	"orgs:write",
	"users:read",
	"users:write",
	"network_elements:read",
	"network_elements:write",
	"policies:read",
	"policies:write",
	"groups:read",
	"groups:write",
	"roles:read",
	"roles:write",
	"egress_routes:read",
	"egress_routes:write",
	"locations:read",
	"audit:read",
	"metrics:read",
	"api_keys:read",
	"api_keys:write",
	"settings:read",
	"settings:write",
	"metaports:read",
	"metaports:write",
	"routing_groups:read",
	"routing_groups:write",
	"peerings:read",
	"peerings:write",
	"metaconnects:read",
	"metaconnects:write",
	"easylinks:read",
	"easylinks:write",
	"tags:read",
	"tags:write",
	"alerts:read",
	"alerts:write",
	"posture_checks:read",
	"posture_checks:write",
	"access_bridges:read",
	"access_bridges:write",
	"trusted_networks:read",
	"trusted_networks:write",
	"certificates:read",
	"certificates:write",
	"apps:read",
	"apps:write",
	"version_controls:read",
	"version_controls:write",
	"access_controls:read",
	"access_controls:write",
	"url_filtering_rules:read",
	"url_filtering_rules:write",
	"threat_categories:read",
	"threat_categories:write",
	"content_categories:read",
	"content_categories:write",
	"pac_files:read",
	"pac_files:write",
	"metaport_clusters:read",
	"metaport_clusters:write",
	"metaport_failovers:read",
	"metaport_failovers:write",
	"cloud_apps:read",
	"cloud_apps:write",
	"file_scanning_rules:read",
	"file_scanning_rules:write",
	"ssl_bypass_rules:read",
	"ssl_bypass_rules:write",
	"enterprise_dns:read",
	"enterprise_dns:write",
	"proxy_access:write",
	"dlp_rules:read",
	"dlp_rules:write",
	"tenant_restrictions:read",
	"tenant_restrictions:write",
}

var nonValidPrivs = []string{
	"abcde",
	"read",
	"write",
	"test123:read",
}

func TestPrivilegesPattern(t *testing.T) {
	for _, priv := range validPrivs {
		assert.True(t, PrivilegesPattern.MatchString(priv))
	}
	for _, priv := range nonValidPrivs {
		assert.False(t, PrivilegesPattern.MatchString(priv))
	}
}

func TestValidateEmail(t *testing.T) {
	cases := map[string]struct {
		Input       string
		ShouldError bool
	}{
		"positive-test": {
			Input:       "test@test.com",
			ShouldError: false,
		},
		"positive-test2": {
			Input:       "a@valid.email",
			ShouldError: false,
		},
		"positive-test-mail-with-tag": {
			Input:       "user.name+tag+sorting@example.com",
			ShouldError: false,
		},
		"positive-test-mail-with-hyphen": {
			Input:       "an.email-with-hyphen@examle.com",
			ShouldError: false,
		},
		"positive-test-mail-with-non-alphabetic": {
			Input:       "#!$%&'*+-/=?^_`{}|~@example.org",
			ShouldError: false,
		},
		"negative-test": {
			Input:       "invalid.email",
			ShouldError: true,
		},
		"negative-test-two-@": {
			Input:       "invalid@e@mail.com",
			ShouldError: true,
		},
		"negative-test-no-domain": {
			Input:       "invalid@email",
			ShouldError: true,
		},
		"negative-test-spaces": {
			Input:       "abc.def @valid.email",
			ShouldError: true,
		},
		"negative-test-backslash": {
			Input:       "abc.def\\\\u00a@valid.email",
			ShouldError: true,
		},
		"negative-test-invalid-local": {
			Input:       "john..doe@example.com",
			ShouldError: true,
		},
		"negative-test-invalid-domain": {
			Input:       "john.doe@example..com",
			ShouldError: true,
		},
		"negative-test-local-starts-with-dot": {
			Input:       ".john.doe@example.com",
			ShouldError: true,
		},
		"negative-test-local-ends-with-dot": {
			Input:       "john.doe.@example.com",
			ShouldError: true,
		},
		"negative-test-domain-starts-with-dot": {
			Input:       "john.doe@.example.com",
			ShouldError: true,
		},
		"negative-test-domain-ends-with-dot": {
			Input:       "john.doe@example.com.",
			ShouldError: true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			diags := ValidateEmail()(tc.Input, nil)
			if diags.HasError() && !tc.ShouldError {
				t.Errorf("%s failed: %+v", name, diags[0])
			}
		})
	}
}

func TestValidateURL(t *testing.T) {
	cases := map[string]struct {
		Input       string
		ShouldError bool
	}{
		"positive-test": {
			Input:       "http://google.com/",
			ShouldError: false,
		},
		"positive-test2": {
			Input:       "https://hooks.slack.com/services/test/1",
			ShouldError: false,
		},
		"positive-test3": {
			Input:       "https://www.dumpsters.com:443",
			ShouldError: false,
		},
		"negative-test-with-invalid-schema": {
			Input:       "http//google.com",
			ShouldError: true,
		},
		"negative-test-without-host-and-schema": {
			Input:       "/foo/bar",
			ShouldError: true,
		},
		"negative-test-host-only": {
			Input:       "google.com",
			ShouldError: true,
		},
		"negative-test-schema-only": {
			Input:       "https",
			ShouldError: true,
		},
		"negative-test-schema-only2": {
			Input:       "https://",
			ShouldError: true,
		},
		"negative-test-empty-string": {
			Input:       "",
			ShouldError: true,
		},
		"negative-test-not-url": {
			Input:       "alskjff#?asf//dfas",
			ShouldError: true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			diags := ValidateURL()(tc.Input, nil)
			if diags.HasError() && !tc.ShouldError {
				t.Errorf("%s failed: %+v", name, diags[0])
			}
		})
	}
}
