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
var hostnameLabelPattern = regexp.MustCompile("^[A-Za-z\\d\\_][A-Za-z\\d\\-\\_]{0,62}[A-Za-z\\d\\_]$")

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

func validateHostName() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		inputString := input.(string)
		if len(inputString) == 0 || len(inputString) > 255 || strings.HasSuffix(inputString, ".") {
			return diag.Errorf("\"%s\" is not a valid hostname", inputString)
		}
		labels := strings.Split(inputString, ".")
		numLabels := len(labels)
		if match, _ := regexp.MatchString("\"[0-9]+$", labels[numLabels-1]); match {
			return diag.Errorf("\"%s\" is not a valid hostname - the TLD must not be all-numeric", inputString)
		}
		for _, l := range labels {
			if !hostnameLabelPattern.MatchString(l) {
				return diag.Errorf("\"%s\" is not a valid hostname", inputString)
			}
		}
		return diags
	}
}

func validateWildcardHostName() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) diag.Diagnostics {
		inputString := input.(string)
		if strings.HasPrefix(inputString, "*.") {
			return validateHostName()(inputString[2:], nil)
		}
		return validateHostName()(inputString, nil)
	}
}

// validateIntRange that integer value is between specified range
func validateIntRange(min, max int) func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		inputInt := input.(int)
		if inputInt < min {
			return diag.Errorf("%d is lower than minimum value %d", inputInt, min)
		}
		if inputInt > max {
			return diag.Errorf("%d is lower than maximum value %d", inputInt, min)
		}
		return diags
	}
}
