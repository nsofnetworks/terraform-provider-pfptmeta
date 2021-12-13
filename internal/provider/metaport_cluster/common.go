package metaport_cluster

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

const (
	description        = "MetaPort cluster defines a group of highly-available MetaPorts that are deployed together in a single data center"
	mappedElementsDesc = "List of mapped element IDs"
	metaportsDesc      = "List of MetaPort IDs"
)

var excludedKeys = []string{"id"}

func metaportClusterRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	m, err := client.GetMetaportCluster(c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
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

func metaportClusterCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	body := client.NewMetaportCluster(d)
	m, err := client.CreateMetaportCluster(c, body)
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

func metaportClusterUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewMetaportCluster(d)
	m, err := client.UpdateMetaportCluster(c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(m, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func metaportClusterDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteMetaportCluster(c, id)
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
