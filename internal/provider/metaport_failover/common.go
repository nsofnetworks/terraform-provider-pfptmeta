package metaport_failover

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

var metaportFailoverExcludedKeys = []string{"id", "failback", "failover"}

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
