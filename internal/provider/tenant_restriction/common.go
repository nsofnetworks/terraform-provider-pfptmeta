package tenant_restriction

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

const (
	description = "Tenant restrictions are used to control access to SaaS cloud applications, such as Office 365." +
		" This is done to prevent unwarranted access to potentially malicious tenants," +
		" or the ones that are forbidden from accessing due to risk of data loss. " +
		"With tenant restrictions, organizations can specify the list of tenants that their users are permitted to access," +
		" blocking all other tenants." +
		" When tenant restrictions are used, the Web Security proxy inserts a list of permitted tenants into" +
		" traffic destined for the relevant cloud app, allowing it to perform the enforcement."
	googleTenantsDesc        = "Configuring this will cause Google to issue security token only for the specified tenants."
	allowServiceAccountDesc  = "Whether to allow access to authenticated service accounts."
	allowConsumerAccessDesc  = "Whether to allow access consumer Google Accounts, such as @gmail.com and @googlemail.com.\n\n"
	allowPersonalDomainsDesc = "Whether to allow Microsoft applications for consumer accounts."
	TenantDirectoryIdDesc    = "The directory ID of the tenant that sets tenant restrictions."
	microsoftTenantsDesc     = "Configuring this will cause Azure AD to issue security token only for the specified tenants. " +
		"Any domain that is registered with a tenant can be used to identify the tenant in this list, as well as the directory ID itself"
)

var ExcludedKeys = []string{"id", "google_config", "microsoft_config"}

func tenantRestrictionToResource(d *schema.ResourceData, tr *client.TenantRestriction) (diags diag.Diagnostics) {
	d.SetId(tr.ID)
	err := client.MapResponseToResource(tr, d, ExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}

	switch tr.Type {
	case "GOOGLE":
		if tr.GoogleConfig != nil {
			googleConfig := []map[string]interface{}{
				{
					"tenants":                tr.GoogleConfig.Tenants,
					"allow_consumer_access":  tr.GoogleConfig.AllowConsumerAccess,
					"allow_service_accounts": tr.GoogleConfig.AllowServiceAccounts,
				},
			}
			err = d.Set("google_config", googleConfig)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	case "MICROSOFT":
		if tr.MicrosoftConfig != nil {
			microsoftConfig := []map[string]interface{}{
				{
					"tenants":                          tr.MicrosoftConfig.Tenants,
					"allow_personal_microsoft_domains": tr.MicrosoftConfig.AllowPersonalMicrosoftDomains,
					"tenant_directory_id":              tr.MicrosoftConfig.TenantDirectoryId,
				},
			}
			err = d.Set("microsoft_config", microsoftConfig)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return
}

func trRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Get("id").(string)
	tr, err := client.GetTenantRestriction(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing tenant restriction %s because it's gone", id)
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	return tenantRestrictionToResource(d, tr)
}
func trCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewTenantRestriction(d)
	tr, err := client.CreateTenantRestriction(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return tenantRestrictionToResource(d, tr)
}

func trUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewTenantRestriction(d)
	tr, err := client.UpdateTenantRestriction(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return tenantRestrictionToResource(d, tr)
}

func trDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteTenantRestriction(ctx, c, id)
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
