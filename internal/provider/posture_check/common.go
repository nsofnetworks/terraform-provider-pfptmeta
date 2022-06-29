package posture_check

import (
	"context"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
)

const (
	description = `Posture checks are administrator-defined sets of criteria allowing to or preventing the user devices from connecting to Proofpoint NaaS.
Administrators can use SQL to create real-world conditions based on underlying information from the operating system and its hardware or choose from a list of common pre-defined conditions.
For example, the presence (or absence) of a specified file or a system process can serve as a pre-condition for letting a device to access the Proofpoint NaaS.
The posture checks can be based on SQL query strings. The queries use the osquery framework, see [osquery.io](https://osquery.io/) for details on osquery.
Posture checks can be viewed and filtered by failure via security logs see [here](https://help.metanetworks.com/knowledgebase/posture_checks) for more details.
`
	enabledDesc = "Defaults to true"
	actionDesc  = "Action to take in case a posture check fails. ENUM: `DISCONNECT`, `NONE`, `WARNING`:\n" +
		"	- **Disconnect** - disconnect device from Proofpoint NaaS.\n" +
		"	- **None** - do nothing, useful during the discovery phase. \n" +
		"   - **Warning** - pop up a warning message, useful during the discovery phase."
	checkDesc           = "Predefined checks. cannot be set with `osquery`."
	minVersionDesc      = "Minimum version required by the check. Required when `type` is `minimum_app_version` or `minimum_os_version`, format: major.minor.patch."
	typeDesc            = "ENUM: `jailbroken_rooted`, `screen_lock_enabled`, `minimum_app_version`, `minimum_os_version`, `malicious_app_detection`, `developer_mode_enabled`."
	whenDesc            = "When to run the check, ENUM: `PRE_CONNECT`, `PERIODIC`."
	applyToOrgDesc      = "Whether to apply to all devices on the organization. Note: this attribute overrides `apply_to_entities`"
	applyToEntitiesDesc = "Entities (users, groups or network elements) to be applied in the posture check."
	exemptEntitiesDesc  = "Entities (users, groups or network elements) which are exempt from the posture check."
	osQueryDesc         = "osquery to use in the posture check, see [here](https://osquery.io/) for more details."
	platformDesc        = "Device platforms that should be applied in the posture check. ENUM: `Android`, `macOS`, `iOS`, `Linux`, `Windows`, `ChromeOS`."
	userMessageDesc     = "Message to be displayed when posture check fails."
	intervalDesc        = "Interval in minutes between checks, mandatory when `when` is set to `PERIODIC`. ENUM: 5, 60."
)

var excludedKeys = []string{"id", "check"}

func postureCheckCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewPostureCheck(d)
	pc, err := client.CreatePostureCheck(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return postureCheckToResource(d, pc)
}
func postureCheckRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)

	pc, err := client.GetPostureCheck(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing posture_check %s because it's gone", id)
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	return postureCheckToResource(d, pc)
}

func postureCheckUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)
	id := d.Id()
	body := client.NewPostureCheck(d)
	pc, err := client.UpdatePostureCheck(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return postureCheckToResource(d, pc)
}
func postureCheckDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeletePostureCheck(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	d.SetId("")
	return diags
}

func postureCheckToResource(d *schema.ResourceData, pc *client.PostureCheck) (diags diag.Diagnostics) {
	d.SetId(pc.ID)
	err := client.MapResponseToResource(pc, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	if pc.Check != nil {
		check := []map[string]interface{}{
			{"min_version": pc.Check.MinVersion, "type": pc.Check.Type},
		}
		err = d.Set("check", check)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return
}
