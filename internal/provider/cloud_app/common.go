package cloud_app

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

const (
	description = `Additional flexibility can be achieved by controlling user access to cloud-based applications.
The administrators can block, redirect to isolation or record any attempt to access applications selected from a 
vast Proofpoint catalog or defined via a specific URL.
`
	appDesc    = "The ID of the [catalog_app](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/data-sources/catalog_app) data-source."
	urlsDesc   = "A list of URLs to associate with this cloud app."
	tenantDesc = `Specific tenant ID of the app on which the cloud application rule should be applied. 
Valid only for catalog apps that have [tenant_corp_id_support](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/catalog_app#tenant_corp_id_support) set to true`
	tenantTypeDesc = "ENUM: `All`, `Personal`, `Corporate` (Defaults to All). " +
		"Valid only for catalog apps that have [tenant_type_support](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/resources/catalog_app#tenant_type_support) set to true"
)

var excludedKeys = []string{"id"}

func cloudAppRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	ca, err := client.GetCloudApp(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing cloud_app %s because it's gone", id)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	err = client.MapResponseToResource(ca, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ca.ID)
	return
}
func cloudAppCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewCloudApp(d)
	ac, err := client.CreateCloudApp(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ac.ID)
	return cloudAppRead(ctx, d, c)
}

func cloudAppUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewCloudApp(d)
	ac, err := client.UpdateCloudApp(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(ac.ID)
	return cloudAppRead(ctx, d, c)
}

func cloudAppDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteCloudApp(ctx, c, id)
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
