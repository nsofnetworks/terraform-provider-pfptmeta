package ssl_bypass_rule

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
	"log"
	"net/http"
)

var excludedKeys = []string{"id"}

const (
	description = "You can define specific content types or domains that are not subject to SSL decryption as their traffic is directed through the proxy. " +
		"The purpose of disabling decryption is to protect personal identification information. " +
		"Currently, the bypassed SSL traffic is not recorded in the Web Security logs."
	applyToOrgDesc = "Indicates whether this SSL bypass rule applies to the org."
	sourcesDesc    = "Users and groups on which the SSL bypass rule should be applied"
	exemptSources  = "Subgroup of 'sources' on which the SSL bypass rule should not be applied"
	priorityDesc   = "Determines the order in which the SSL bypass rules are evaluated. " +
		"The order is significant since among all the SSL bypass rules that are relevant to a specific user, " +
		"the one with the highest priority (smaller priority value) is the one to determine the SSL Bypass enforcement applied to that user."
	uncategorizedUelDesc = "Whether to SSL bypass uncategorized URLs."
	contentTypesDesc     = "A List of content types. If a domain is found to be categorized under at least of one of them, it will be bypassed. " + common.ContentTypesDoc
	domainsDesc          = "A list of domains to SSL bypass."
	actionDesc           = "Enum: `BYPASS`, `INTERCEPT`.\n" +
		"The action to take in case of a match"
)

func parseSslBypassRule(d *schema.ResourceData, pf *client.SSLBypassRule) diag.Diagnostics {
	d.SetId(pf.ID)
	err := client.MapResponseToResource(pf, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diag.Diagnostics{}
}

func sslBypassRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	pf, err := client.GetSSLBypassRule(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing ssl bypass rule %s because it's gone", id)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	return parseSslBypassRule(d, pf)
}

func pacFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewSSLBypassRule(d)
	pf, err := client.CreateSSLBypassRule(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return parseSslBypassRule(d, pf)
}

func pacFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewSSLBypassRule(d)
	pf, err := client.UpdateSSLBypassRule(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return parseSslBypassRule(d, pf)
}

func pacFileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteSSLBypassRule(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return
}
