package common

import (
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var numericPattern = regexp.MustCompile("^[0-9]{1,30}$")
var alphabetPattern = regexp.MustCompile("^[a-zA-Z0-9]{1,30}$")
var hostnameLabelPattern = regexp.MustCompile("^[A-Za-z\\d\\_][A-Za-z\\d\\-\\_]{0,62}[A-Za-z\\d\\_]$")
var TagPattern = regexp.MustCompile("^[a-zA-Z0-9-_]+$")
var PrivilegesPattern = regexp.MustCompile("^[a-z_]+:(read|write)$")
var HttpHeaderPattern = regexp.MustCompile("^([\\w\\-]+):(.*)$")
var DomainPattern = regexp.MustCompile("^(?:[a-z0-9](?:[a-z0-9-_]{0,61}[a-z0-9])?\\.)+[a-z0-9][a-z0-9-_]{0,61}[a-z]$")
var AccessIdPattern = regexp.MustCompile("^([A-Za-z0-9_-]={0,2}){40,50}$")

func containsString(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}

func containsInt(v int, a []int) bool {
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

// ValidateIntENUM validates ENUM values
func ValidateIntENUM(enum ...int) func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		switch {
		case containsInt(input.(int), enum):
			return diags
		default:
			return diag.Errorf("%d is not one of %+v", input, enum)
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

// ValidateStringToIntRange that integer value given as string is between specified range
func ValidateStringToIntRange(min, max int) func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) diag.Diagnostics {
		var diags diag.Diagnostics
		origInput := input.(string)
		if origInput == "" {
			return diags
		}
		inputInt, err := strconv.Atoi(input.(string))
		if err != nil {
			return diag.FromErr(err)
		}
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

func ValidateCustomUrlOrIPV4() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) (diags diag.Diagnostics) {
		inputString := input.(string)
		if strings.HasPrefix(inputString, ".") {
			return ValidateHostName()(inputString[1:], nil)
		}
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

func ValidateHTTPNetLocation() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, _ cty.Path) (diags diag.Diagnostics) {
		inputString := input.(string)
		h, err := url.ParseRequestURI(inputString)
		if err != nil {
			return diag.Errorf(
				"\"%s\" is not a valid host: %s", inputString, err)
		}
		if h.Scheme != "http" && h.Scheme != "https" {
			return diag.Errorf(
				"\"%s\" is not a valid host: should have http or https schema only, got \"%s\"",
				inputString, h.Scheme)
		}
		if h.Path != "" {
			return diag.Errorf(
				"\"%s\" is not a valid host: path is not allowed - got \"%s\"",
				inputString, h.Path)
		}
		if h.RawQuery != "" {
			return diag.Errorf(
				"\"%s\" is not a valid host: query params are not allowed - got \"%s\"",
				inputString, h.RawQuery)
		}
		if h.RawFragment != "" {
			return diag.Errorf(
				"\"%s\" is not a valid host: fragment is not allowed - got \"%s\"",
				inputString, h.RawFragment)
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

// ComposeOrValidations will take a list of validation functions and will allow the validated attribute
// if at least one function has approved the attribute
func ComposeOrValidations(fs ...func(interface{}, cty.Path) diag.Diagnostics) func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, p cty.Path) (diags diag.Diagnostics) {
		for _, f := range fs {
			fDiags := f(input, p)
			if !fDiags.HasError() {
				return diag.Diagnostics{}
			}
			diags = append(diags, fDiags...)
		}
		return
	}
}

func ValidateDNS() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, path cty.Path) (diags diag.Diagnostics) {
		inputString := input.(string)
		if strings.Contains(inputString, "_") || !strings.Contains(inputString, ".") {
			return diag.Errorf("\"%s\" is not a valid domain name", inputString)
		}
		if ValidateHostName()(inputString, path).HasError() {
			return diag.Errorf("\"%s\" is not a valid domain name", inputString)
		}
		return
	}
}

func ValidateDomainName() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, path cty.Path) (diags diag.Diagnostics) {
		inputString := input.(string)
		if inputString == "" {
			return
		}
		if ValidateHostName()(inputString, path).HasError() {
			return diag.Errorf("\"%s\" is not a valid domain name", inputString)
		}
		return
	}
}

func ValidateCIDR4() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, path cty.Path) (diags diag.Diagnostics) {
		inputString := input.(string)
		ipv4addr, ipv4Net, err := net.ParseCIDR(inputString)
		if err != nil {
			return diag.FromErr(err)
		}
		if len(ipv4Net.Mask) != net.IPv4len || !ipv4addr.Equal(ipv4Net.IP) {
			return diag.Errorf("\"%s\" is not a valid IPV4-CIDR", inputString)
		}
		return
	}
}

func ValidatePEMCert() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, path cty.Path) (diags diag.Diagnostics) {
		inputString := input.(string)
		block, _ := pem.Decode([]byte(inputString))
		if block == nil || block.Type != "CERTIFICATE" {
			return diag.Errorf("failed to decode PEM block containing certificate")
		}
		_, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return diag.FromErr(err)
		}
		return
	}
}

func ValidateIsoTimeFormat() func(interface{}, cty.Path) diag.Diagnostics {
	return func(input interface{}, path cty.Path) (diags diag.Diagnostics) {
		inputString := input.(string)
		_, err := time.Parse(time.RFC3339, inputString)
		if err != nil {
			return diag.FromErr(err)
		}
		return
	}
}
