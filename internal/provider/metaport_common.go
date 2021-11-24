package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

var metaportExcludedKeys = []string{"id"}
var metaportFailoverExcludedKeys = []string{"id", "failback", "failover"}

func metaportRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	m, err := client.GetMetaport(c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	err = client.MapResponseToResource(m, d, metaportExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(m.ID)
	return diags
}

func metaportCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	body := client.NewMetaport(d)
	m, err := client.CreateMetaport(c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(m.ID)
	err = client.MapResponseToResource(m, d, metaportExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func metaportUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewMetaport(d)
	m, err := client.UpdateMetaport(c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	err = client.MapResponseToResource(m, d, metaportExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	return diags
}

func metaportDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteMetaport(c, id)
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
	err = client.MapResponseToResource(m, d, metaportExcludedKeys)
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
	err = client.MapResponseToResource(m, d, metaportExcludedKeys)
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
	err = client.MapResponseToResource(m, d, metaportExcludedKeys)
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

func metaportFailoverRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	m, err := client.GetMetaportFailover(c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
			return diags
		} else {
			return diag.FromErr(err)
		}
	}
	return metaportFailoverToResource(d, m)
}

func metaportFailoverCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewMetaportFailover(d)
	m, err := client.CreateMetaportFailover(c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return metaportFailoverToResource(d, m)
}

func metaportFailoverUpdate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewMetaportFailover(d)
	m, err := client.UpdateMetaportFailover(c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return metaportFailoverToResource(d, m)
}

func metaportFailoverDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteMetaportFailover(c, id)
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

func metaportFailoverToResource(d *schema.ResourceData, m *client.MetaportFailover) diag.Diagnostics {
	var diags diag.Diagnostics
	d.SetId(m.ID)
	err := client.MapResponseToResource(m, d, metaportFailoverExcludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	if m.FailOver != nil {
		failoverToResource := []map[string]interface{}{
			{"delay": m.FailOver.Delay, "threshold": m.FailOver.Threshold, "trigger": m.FailOver.Trigger},
		}
		err = d.Set("failover", failoverToResource)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if m.FailBack != nil {
		failbackToReource := []map[string]interface{}{
			{"trigger": m.FailBack.Trigger},
		}
		err = d.Set("failback", failbackToReource)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	return diags
}
