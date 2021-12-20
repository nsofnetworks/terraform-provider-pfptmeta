package common

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strings"
)

var numericPattern = regexp.MustCompile("^[0-9]{1,30}$")
var alphabetPattern = regexp.MustCompile("^[a-zA-Z0-9]{1,30}$")
var hostnameLabelPattern = regexp.MustCompile("^[A-Za-z\\d\\_][A-Za-z\\d\\-\\_]{0,62}[A-Za-z\\d\\_]$")
var TagPattern = regexp.MustCompile("^[a-zA-Z0-9-_]+$")
var PrivilegesPattern = regexp.MustCompile("^[a-z_]+:(read|write)$")
var HttpHeaderPattern = regexp.MustCompile("^([\\w\\-]+):(.*)$")

func containsString(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}

// ValidateStringENUM validates ENUM values
func ValidateStringENUM(enum ...string) func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		inputString := input.(string)
		switch {
		case containsString(inputString, enum):
			return diags
		default:
			return append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  fmt.Sprintf("\"%s\" is not one of %+v", inputString, enum),
			})
		}
	}
}

func ValidateID(numeric bool, prefixes ...string) func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		ID := input.(string)
		parts := strings.Split(ID, "-")
		if len(parts) != 2 {
			return diag.Errorf("\"%s\" should be of the form <prefix>-<unique>", ID)
		}
		switch {
		case !containsString(parts[0], prefixes):
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

func ValidatePattern(pattern *regexp.Regexp) func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		inputString := input.(string)
		if !pattern.MatchString(inputString) {
			return diag.Errorf("\"%s\" does not match pattern \"%s\"", inputString, pattern.String())
		}
		return diags
	}
}

func ValidateHostName() func(interface{}, cty.Path) diag.Diagnostics {
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

func ValidateWildcardHostName() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) diag.Diagnostics {
		inputString := input.(string)
		if strings.HasPrefix(inputString, "*.") {
			return ValidateHostName()(inputString[2:], nil)
		}
		return ValidateHostName()(inputString, nil)
	}
}

// ValidateIntRange that integer value is between specified range
func ValidateIntRange(min, max int) func(interface{}, cty.Path) diag.Diagnostics {
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

func ValidateHostnameOrIPV4() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) (diags diag.Diagnostics) {
		inputString := input.(string)
		hostnameValidation := ValidateHostName()(inputString, nil)
		ipv4Validation := net.ParseIP(inputString) != nil
		if hostnameValidation.HasError() && !ipv4Validation {
			return diag.Errorf("\"%s\" is not a valid hostname or ipv4", inputString)
		}
		return
	}
}

func ValidateEmail() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) (diags diag.Diagnostics) {
		inputString := input.(string)
		if len(inputString) > 254 {
			return diag.Errorf("\"%s\" is not a valid email - cannot be longer than 254 characters", inputString)
		}
		_, err := mail.ParseAddress(inputString)
		if err != nil {
			return diag.Errorf("\"%s\" is not a valid email", inputString)
		}
		return
	}
}

func ValidateURL() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) (diags diag.Diagnostics) {
		inputString := input.(string)
		_, err := url.ParseRequestURI(inputString)
		if err != nil {
			return diag.Errorf("\"%s\" is not a valid url %s", inputString, err)
		}
		return
	}
}

func ValidateJson() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) (diags diag.Diagnostics) {
		inputString := input.(string)
		var js json.RawMessage
		err := json.Unmarshal([]byte(inputString), &js)
		if err != nil {
			return diag.Errorf("\"%.200s\" is not a valid json. %s", inputString, err)
		}
		return
	}
}
