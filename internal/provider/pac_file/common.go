package pac_file

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

var excludedKeys = []string{"id"}

const (
	description = "Web traffic inspection is further enhanced by means of traffic steering rules (Implemented as Proxy Auto Config file), " +
		"installed by the web security engine after the user has been onboarded. " +
		"It is a JavaScript-based file that uses logical statements to determine which traffic is routed through the proxy and which traffic bypasses it. " +
		"Each tenant has a default rule supplied by Proofpoint, based on the best practice recommendations. " +
		"However, the administrators can decide to override the default rule with customized traffic steering rules that are better suited for their organization. " +
		"Once created, the traffic steering rule is distributed to intended users to be hosted locally on their machines. Afterwards, the rule can be updated at any time."
	applyToOrgDesc = "Indicates whether this PAC file applies to the org."
	sourcesDesc    = "Users and groups on which the PAC file should be applied."
	exemptSources  = "Subgroup of `sources` on which the PAC file should not be applied."
	priorityDesc   = "Determines the order in which the PAC files are being matched. Lower priority value means the PAC file will be matched earlier."
	hasContentDesc = "Whether the PAC file object has content associated with it."
	contentDesc    = "The content of the PAC file"
)

func pacFileToResource(ctx context.Context, d *schema.ResourceData, c *client.Client, pf *client.PacFile) diag.Diagnostics {
	d.SetId(pf.ID)
	err := client.MapResponseToResource(pf, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	if pf.HasContent {
		content, err := client.GetPacFileContent(ctx, c, pf.ID)
		if err != nil {
			return diag.FromErr(err)
		}
		err = d.Set("content", content)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return diag.Diagnostics{}
}

func pacFileRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	pf, err := client.GetPacFile(ctx, c, id)
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
	if err != nil {
		return diag.FromErr(err)
	}
	diags = parsePacFile(ctx, d, c, pf)
	if diags.HasError() {
		return diags
	}
	if content := d.Get("content").(string); content != "" {
		err = client.PutPacFileContent(ctx, c, pf.ID, content)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return pacFileToResource(ctx, d, c, pf)
}

func pacFileUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewPacFile(d)
	pf, err := client.UpdatePacFile(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	diags = parsePacFile(ctx, d, c, pf)
	if diags.HasError() {
		return diags
	}
	if _, content := d.GetChange("content"); content.(string) != "" {
		err = client.PutPacFileContent(ctx, c, pf.ID, content.(string))
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
