package certificate

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

const (
	description = "SSL certificate. It is used mostly to allow EasyLinks to utilize HTTPS, when operating with the `redirect` or `native` access types."
	sansDesc    = "List of certificate SANs"
	stateDesc   = "Certificate state, can be one of the following:\n" +
		"	- **Pending** - Initial state that may take several minutes. During this stage, a request is sent to certification authority and the system is waiting for the certificate approval.\n" +
		"	- **OK** - Certificate has been validated by the certification authority and ready for use.\n" +
		"	- **Warning** - Certificate is valid, but it is to expire within 30 days. DNS check attempts for the certificate renewal have failed.\n" +
		"	- **Error** - Certificate has expired, all DNS checks have failed so far, and no renewal attempts are being made.\n"
)

var excludedKeys = []string{"id"}

func certificateRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	cert, err := client.GetCertificate(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing certificate %s because it's gone", id)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	err = client.MapResponseToResource(cert, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cert.ID)
	return
}

func certificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewCertificate(d)
	cert, err := client.CreateCertificate(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(cert.ID)
	err = client.MapResponseToResource(cert, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func certificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewCertificate(d)
	cert, err := client.UpdateCertificate(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(cert, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func certificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteCertificate(ctx, c, id)
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
