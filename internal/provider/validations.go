package provider

import (
	"fmt"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func contains(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}

// validateENUM validates ENUM values
func validateENUM(enum... string) func(interface{}, cty.Path) diag.Diagnostics {
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