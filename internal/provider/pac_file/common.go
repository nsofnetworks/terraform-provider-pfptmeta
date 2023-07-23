package pac_file

import (
	"context"
	"fmt" // NADAV REMOVE
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

var excludedKeys = []string{"id"}

const pacTypeManaged string = "managed"
const pacTypeBringYourOwn string = "bring_your_own"

const (
	description = "Web traffic inspection is further enhanced by means of traffic steering rules (Implemented as Proxy Auto Config file), " +
		"installed by the web security engine after the user has been onboarded. " +
		"It is a JavaScript-based file that uses logical statements to determine which traffic is routed through the proxy and which traffic bypasses it. " +
		"Each tenant has a default rule supplied by Proofpoint, based on the best practice recommendations. " +
		"However, the administrators can decide to override the default rule with customized traffic steering rules that are better suited for their organization. " +
		"Once created, the traffic steering rule is distributed to intended users to be hosted locally on their machines. Afterwards, the rule can be updated at any time."
	applyToOrgDesc     = "Indicates whether this PAC file applies to the org."
	sourcesDesc        = "Users and groups on which the PAC file should be applied."
	exemptSources      = "Subgroup of `sources` on which the PAC file should not be applied."
	priorityDesc       = "Determines the order in which the PAC files are being matched. Lower priority value means the PAC file will be matched earlier."
	hasContentDesc     = "Whether the PAC file object has content associated with it."
	contentDesc        = "The content of the PAC file"
	typeDesc           = "Indicates whether this PAC file has '" + pacTypeManaged + "' or '" + pacTypeBringYourOwn + "' content type"
	managedContentDesc = "Lists of domains, Cloud App ids and IP Network ids which will automatically " +
		"be monitored for changes and will update the raw content of the PAC file according " +
		"to changes in these Cloud Apps and IP Network objects."
	managedContentDomainsDesc   = "domains to be used as is"
	managedContentCloudAppsDesc = "Ids of Cloud Apps to monitor for changes. The domains from these " +
		"will be added (and updated) to the raw content of the PAC file"
	managedContentIPNetworksDesc = "Ids of Cloud Apps to monitor for changes. The domains from these " +
		"will be added (and updated) to the raw content of the PAC file"
)

func pacFileToResource(ctx context.Context, d *schema.ResourceData, c *client.Client, pf *client.PacFile) diag.Diagnostics {
	d.SetId(pf.ID)
	fmt.Println("[NADAV] pacFileToResource [0]")
	err := client.MapResponseToResource(pf, d, excludedKeys)
	if err != nil {
		fmt.Println("[NADAV] pacFileToResource [0.err]")
		return diag.FromErr(err)
	}
	fmt.Println("[NADAV] pacFileToResource [1]")
	if pf.HasContent {
		content, err := client.GetPacFileContent(ctx, c, pf.ID)
		fmt.Println("[NADAV] pacFileToResource [1.1]")
		if err != nil {
			fmt.Println("[NADAV] pacFileToResource [1.1.err]")
			return diag.FromErr(err)
		}
		err = d.Set("content", content)
		fmt.Println("[NADAV] pacFileToResource [1.2]")
		if err != nil {
			fmt.Println("[NADAV] pacFileToResource [1.2.err]")
			return diag.FromErr(err)
		}
	}
	fmt.Printf("[NADAV] pacFileToResource [2]: pf: %+v\n", *pf)
	if pf.Type == pacTypeManaged {
		managed_content_ptr, err := client.GetPacFileManagedContent(ctx, c, pf.ID)
		fmt.Printf("[NADAV] pacFileToResource [2.1]: managed_content_ptr: %+v\n", *managed_content_ptr)
		if err != nil {
			fmt.Println("[NADAV] pacFileToResource [2.1.err]")
			return diag.FromErr(err)
		}
		managed_content_lst := []map[string]interface{}{}
		fmt.Println("[NADAV] pacFileToResource [2.2]")
		if managed_content_ptr != nil {
			fmt.Println("[NADAV] pacFileToResource [2.3]")
			content_map := make(map[string]interface{})
			if len(managed_content_ptr.Domains) > 0 {
				content_map["domains"] = managed_content_ptr.Domains
			} else {
				content_map["domains"] = []string{}
			}
			if len(managed_content_ptr.CloudApps) > 0 {
				content_map["cloud_apps"] = managed_content_ptr.CloudApps
			} else {
				content_map["cloud_apps"] = []string{}
			}
			if len(managed_content_ptr.IpNetworks) > 0 {
				content_map["ip_networks"] = managed_content_ptr.IpNetworks
			} else {
				content_map["ip_networks"] = []string{}
			}
			managed_content_lst = append(managed_content_lst, content_map)
		}
		fmt.Printf("[NADAV] pacFileToResource [2.4]: managed_content_lst: %+v\n", managed_content_lst)
		err = d.Set("managed_content", managed_content_lst)
		if err != nil {
			fmt.Println("[NADAV] pacFileToResource [2.4.err]")
			return diag.FromErr(err)
		}
		fmt.Printf("[NADAV] pacFileToResource [2.5]: d: %+v", *d)
	}
	return diag.Diagnostics{}
}

func pacFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	pf, err := client.GetPacFile(ctx, c, id)
	fmt.Println("[NADAV] pacFileRead")
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing pac file %s because it's gone", id)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	return pacFileToResource(ctx, d, c, pf)
}

func pacFileCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	body := client.NewPacFile(d)
	pf, err := client.CreatePacFile(ctx, c, body)
	fmt.Println("[NADAV] pacFileCreate")
	if err != nil {
		return diag.FromErr(err)
	}
	if content := d.Get("content").(string); content != "" {
		if d.Get("type") != pacTypeBringYourOwn {
			return diag.Errorf("Content cannot be set for a bring_your_own PAC type")
		}
		err = client.PutPacFileContent(ctx, c, pf.ID, content)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if managed_content_rsrc := d.Get("managed_content"); len(managed_content_rsrc.([]interface{})) != 0 {
		if d.Get("type") != pacTypeManaged {
			return diag.Errorf("Managed Content can only be set for a managed PAC type")
		}
		managed_content := client.NewManagedContent(d)
		err = client.PatchPacFileManagedContent(ctx, c, pf.ID, managed_content)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return pacFileToResource(ctx, d, c, pf)
}

func pacFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	if d.HasChange("type") {
		return diag.Errorf("Cannot change type of existing PAC")
	}
	body := client.ModifiedPacFile(d)
	pf, err := client.UpdatePacFile(ctx, c, id, body)
	fmt.Println("[NADAV] pacFileUpdate [0]")
	if err != nil {
		return diag.FromErr(err)
	}
	if d.HasChange("content") {
		if d.Get("type") != pacTypeBringYourOwn {
			return diag.Errorf("Content can only be update for bring_your_own PAC type")
		}
		_, new_content := d.GetChange("content")
		err = client.PutPacFileContent(ctx, c, pf.ID, new_content.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	}
	fmt.Println("[NADAV] pacFileUpdate [1]")
	if d.HasChange("managed_content") {
		fmt.Println("[NADAV] pacFileUpdate [1.1]")
		if d.Get("type") != pacTypeManaged {
			return diag.Errorf("Managed Content can only be update for managed PAC type")
		}
		new_content := client.NewManagedContent(d)
		fmt.Printf("[NADAV] pacFileUpdate [1.2]: new_content: %+v\n", *new_content)
		err = client.PatchPacFileManagedContent(ctx, c, pf.ID, new_content)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return pacFileToResource(ctx, d, c, pf)
}

func pacFileDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeletePacFile(ctx, c, id)
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
