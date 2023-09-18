package aac_rule

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/provider/common"
	"log"
	"net/http"
)

var excludedKeys = []string{"id", "suspicious_login"}

const (
	description = "Adaptive access control rule for protecting users connecting to service provider application " +
		"under risky conditions"
	priorityDesc = "Determines the order in which the aac rules are being matched. " +
		"Lower priority indicates that the AAC rule is matched earlier"
	actionDesc       = "The action to enforce when rule is matched to a connection"
	appIdsDesc       = "IDs of the apps that the AAC rule is applied to"
	applyAllAppsDesc = "Indicates whether this rule applies to all apps of the org, regardless whether such " +
		"apps are specified in app_ids. Note: this attribute overrides app_ids"
	sourcesDesc              = "Users and groups that the rule is applied to"
	exemptSources            = "Subgroup of 'sources' to which the AAC rule is not applied"
	suspiciousLoginDesc      = "Determines if the rule applies at suspicious or non-suspicious login. Options: any, suspicious, safe"
	expressionDesc           = "Defines filtering expressions to to provide user granularity in AAC rule application"
	networksDesc             = "List of IP network IDs that the rule is applied to"
	locationsDesc            = "List of locations that the rule is applied to. Each country is represented by an Alpha-2 code (ISO-3166). Enum: " + common.CountriesDoc
	IPDeputationsDesc        = "List of IP reputations that the rule is applied to"
	CertificateIdDesc        = "The root/intermediate certificate ID of managed devices that the rule is applied to"
	notificationChannelsDesc = "List of notification channel IDs"
)

func aacRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	a, err := client.GetAacRule(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing aac rule %s because it's gone", id)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId(a.ID)
	err = client.MapResponseToResource(a, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("suspicious_login", client.ParseAacSuspiciousLoginBoolToStr(a))
	return
}

func aacRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewAacRule(d)
	a, err := client.CreateAacRule(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(a.ID)
	err = client.MapResponseToResource(a, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("suspicious_login", client.ParseAacSuspiciousLoginBoolToStr(a))
	return
}

func aacRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewAacRule(d)
	a, err := client.UpdateAacRule(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(a, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.Set("suspicious_login", client.ParseAacSuspiciousLoginBoolToStr(a))
	return
}

func aacRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteAacRule(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	return
}
