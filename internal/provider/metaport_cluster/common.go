package metaport_cluster

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
	"strings"
)

const (
	description        = "MetaPort cluster defines a group of highly-available MetaPorts that are deployed together in a single data center"
	mappedElementsDesc = "List of mapped element IDs"
	metaportsDesc      = "List of MetaPort IDs"
)

var excludedKeys = []string{"id"}

func metaportClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)
	var err error
	var m *client.MetaportCluster
	id, exists := d.GetOk("id")
	if exists {
		m, err = client.GetMetaportCluster(ctx, c, id.(string))
	}
	name, exists := d.GetOk("name")
	if exists {
		m, err = client.GetMetaportClustertByName(ctx, c, name.(string))
	}
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing metaport cluster %s because it's gone", id)
			d.SetId("")
			return diags
		} else if strings.HasPrefix(err.Error(), "could not find metaport") {
			log.Printf("[WARN] Removing metaport cluster %s because it's gone", name)
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

func metaportClusterCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	body := client.NewMetaportCluster(d)
	m, err := client.CreateMetaportCluster(ctx, c, body)
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

func metaportClusterUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewMetaportCluster(d)
	m, err := client.UpdateMetaportCluster(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(m, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func metaportClusterDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteMetaportCluster(ctx, c, id)
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
