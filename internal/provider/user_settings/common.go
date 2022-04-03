package user_settings

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
	"strconv"
)

var excludedKeys = []string{"id", "max_devices_per_user"}

const (
	description = "The `pfptmeta_user_settings` resource is a tool with which the administrator can configure specific user settings for particular groups.\n" +
		"For example, an organization’s security policy may require that a specific contractor’s group is prompted for re-authentication after x minutes," +
		" or that this group of users can only log in from a single device or x number of devices.\n" +
		"In addition, the administrator can choose the type of authentication factor that should be applied or if users can only log in using SSO."
	applyOnOrgDesc        = "Indicates whether this user setting applies to the entire org. Note: this attribute overrides `apply_to_entities`."
	applyToEntitiesDesc   = "Entities (users, groups or network elements) that the user settings will be applied to."
	maxDevicesPerUserDesc = "Integer wrapped as string. Provides the administrator the flexibility to restrict how many devices the user can own or authenticate from."
	mfaRequiredDesc       = "Forces the user for second factor authentication when logging in to Proofpoint NaaS. Enabling this enforces the user to authenticate also by a second factor, as specified by `allowed_factors` parameter."
	allowedFactorsDesc    = "When users are configured to authenticate locally with MFA, you can choose which second authentication factors will be visible to this user group." +
		" The allowed values are: `SMS`, `SOFTWARE_TOTP`, `VOICECALL`, `EMAIL`.\n" +
		"This applies ONLY to local Proofpoint accounts, not to accounts that authenticate via external IdPs (SSO)."
	passwordExpirationDesc = "Allows the administrator to set how often (in days) the end user should set a new login password."
	prohibitedOsDesc       = "Allows the administrator to select operating systems which are prohibited from onboarding. ENUM: `Android`, `macOS`, `iOS`, `Linux`, `Windows`, `ChromeOS`"
	proxyPopsDesc          = "Type of proxy_pops the user will use:\n" +
		"	- **ALL_POPS** - connect to the nearest Point-of-Presence regardless to whether this PoP was upgraded for static IP use or not." +
		"	- **POPS_WITH_DEDICATED_IPS** - enable the use of PoPs with dedicated IP ranges provided by Proofpoint."
	ssoMandatoryDesc = "Force the user into SSO authentication, via the configured IdP." +
		" If this option is enabled and the user attempts to login without SSO, the following message is displayed: *Login without SSO is not allowed by system administrator*."
)

func userSettingsToResource(d *schema.ResourceData, us *client.UserSettings) (diags diag.Diagnostics) {
	d.SetId(us.ID)
	err := client.MapResponseToResource(us, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	if us.MaxDevicesPerUser == nil {
		err = d.Set("max_devices_per_user", nil)
	} else {
		err = d.Set("max_devices_per_user", strconv.Itoa(*us.MaxDevicesPerUser))
	}
	if err != nil {
		return diag.FromErr(err)
	}
	return
}
func userSettingsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	us, err := client.GetUserSettings(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	return userSettingsToResource(d, us)
}
func userSettingsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewUserSettings(d)
	us, err := client.CreateUserSettings(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return userSettingsToResource(d, us)
}

func userSettingsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewUserSettings(d)
	us, err := client.UpdateUserSettings(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return userSettingsToResource(d, us)
}

func userSettingsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteUserSettings(ctx, c, id)
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
