package provider

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"regexp"
	"strings"
)

var numericPattern = regexp.MustCompile("^[0-9]{1,30}$")
var alphabetPattern = regexp.MustCompile("^[a-zA-Z0-9]{1,30}$")

func contains(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}

// validateENUM validates ENUM values
func validateENUM(enum ...string) func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		inputString := input.(string)
		switch {
		case contains(inputString, enum):
			return diags
		default:
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("\"%s\" is not one of %+v", inputString, enum),
			})
		}
	}
}

func validateID(numeric bool, prefixes ...string) func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		ID := input.(string)
		parts := strings.Split(ID, "-")
		if len(parts) != 2 {
			return diag.Errorf("\"%s\" should be of the form <prefix>-<unique>", ID)
		}
		switch {
		case !contains(parts[0], prefixes):
			return diag.Errorf("\"%s\" should have a prefix of %v", ID, prefixes)
		}
		if numeric {
			if !numericPattern.MatchString(parts[1]) {
				return diag.Errorf("\"%s\" should have a numeric suffix", ID)
			}
		} else {
			if !alphabetPattern.MatchString(parts[1]) {
				return diag.Errorf("\"%s\" should have a alphabet suffix", ID)
			}
		}
		return diags
	}
}

func validatePattern(pattern *regexp.Regexp) func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		inputString := input.(string)
		if !pattern.MatchString(inputString) {
			return diag.Errorf("\"%s\" does not match pattern \"%s\"", inputString, pattern.String())
		}
		return diags
	}
}
