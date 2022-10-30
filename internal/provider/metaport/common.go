package metaport

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
	"strings"
)

var excludedKeys = []string{"id"}

const (
	description = "MetaPort is a lightweight virtual appliance that enables the secure authenticated interface " +
		"interact between existing servers and the Proofpoint NaaS cloud. " +
		"Once configured, metaports enable users to access your applications via the Proofpoint cloud."
	mappedElementsDesc       = "List of mapped element IDs"
	notificationChannelsDesc = "List of notification channel IDs"
)

func metaportRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	var err error
	var m *client.Metaport
	id, exists := d.GetOk("id")
	if exists {
		m, err = client.GetMetaport(ctx, c, id.(string))
	}
	name, exists := d.GetOk("name")
	if exists && m == nil {
		m, err = client.GetMetaportByName(ctx, c, name.(string))
	}
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing metaport %s because it's gone", id)
			d.SetId("")
			return diags
		} else if strings.HasPrefix(err.Error(), "could not find metaport") {
			log.Printf("[WARN] Removing metaport %s because it's gone", name)
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	err = client.MapResponseToResource(m, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(m.ID)
	return diags
}

func metaportCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	body := client.NewMetaport(d)
	m, err := client.CreateMetaport(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(m.ID)
	err = client.MapResponseToResource(m, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func metaportUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewMetaport(d)
	m, err := client.UpdateMetaport(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(m, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func metaportDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteMetaport(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		} else {
			return diag.FromErr(err)
		}
	}
	return diags
}
