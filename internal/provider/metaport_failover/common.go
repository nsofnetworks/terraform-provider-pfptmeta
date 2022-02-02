package metaport_failover

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"net/http"
)

const (
	description        = "MetaPort failover defines a failover model between a primary and a secondary MetaPort clusters."
	mappedElementsDesc = "List of mapped element IDs"
	cluster1Desc       = "Priority #1 MetaPort cluster ID. This cluster is active by default. " +
		"When failover condition is met for this cluster, the higher priority cluster becomes active."
	cluster2Desc             = "Priority #2 MetaPort cluster ID. This cluster becomes active, when failover condition is met for a lower priority cluster."
	failbackDesc             = "Primary to secondary cluster switchover."
	failoverDesc             = "Secondary to primary cluster switchover."
	failoverDelayDesc        = "Number of minutes to wait before execution of failover, defaults to 1."
	failoverThresholdDesc    = "Minimum number of healthy MetaPorts to keep a cluster active. Zero (0) denotes all MetaPorts in a cluster."
	triggerDesc              = "ENUM: [auto, manual], defaults to auto."
	notificationChannelsDesc = "List of notification channel IDs"
)

var excludedKeys = []string{"id", "failback", "failover"}

func metaportFailoverRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	m, err := client.GetMetaportFailover(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			d.SetId("")
		}
		return diag.FromErr(err)
	}
	return metaportFailoverToResource(d, m)
}

func metaportFailoverCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	body := client.NewMetaportFailover(d)
	m, err := client.CreateMetaportFailover(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return metaportFailoverToResource(d, m)
}

func metaportFailoverUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewMetaportFailover(d)
	m, err := client.UpdateMetaportFailover(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return metaportFailoverToResource(d, m)
}

func metaportFailoverDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := meta.(*client.Client)

	id := d.Id()
	_, err := client.DeleteMetaportFailover(ctx, c, id)
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
	err := client.MapResponseToResource(m, d, excludedKeys)
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
