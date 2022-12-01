package time_frame

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"log"
	"net/http"
)

var excludedKeys = []string{"id", "org_id", "start_time", "end_time"}

const (
	description = "You can define customized periods of time to be used as time frames for other resources. " +
		"If they are in use, the rule takes effect only within the defined period and is not enforced outside this window."
	daysDesc    = "ENUM: `sunday`,`monday`,`tuesday`,`wednesday`,`thursday`,`friday`,`saturday`"
	endTimeDesc = "When `end_time` <= `start_time` the time frame will extend to the next day."
)

func timeFrameToResource(d *schema.ResourceData, tf *client.TimeFrame) (diags diag.Diagnostics) {
	d.SetId(tf.ID)
	err := client.MapResponseToResource(tf, d, excludedKeys)
	if err != nil {
		return diag.FromErr(err)
	}
	startTime := []map[string]interface{}{
		{
			"hour":   tf.StartTime.Hour,
			"minute": tf.StartTime.Minute,
		},
	}
	err = d.Set("start_time", startTime)

	endTime := []map[string]interface{}{
		{
			"hour":   tf.EndTime.Hour,
			"minute": tf.EndTime.Minute,
		},
	}
	err = d.Set("end_time", endTime)
	if err != nil {
		return diag.FromErr(err)
	}
	return
}

func timeFrameRead(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	id := d.Get("id").(string)
	c := meta.(*client.Client)
	tf, err := client.GetTimeFrame(ctx, c, id)
	if err != nil {
		errResponse, ok := err.(*client.ErrorResponse)
		if ok && errResponse.Status == http.StatusNotFound {
			log.Printf("[WARN] Removing time frame %s because it's gone", id)
			d.SetId("")
			return
		} else {
			return diag.FromErr(err)
		}
	}
	return timeFrameToResource(d, tf)
}
func timeFrameCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	body := client.NewTimeFrame(d)
	tf, err := client.CreateTimeFrame(ctx, c, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return timeFrameToResource(d, tf)
}

func timeFrameUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)

	id := d.Id()
	body := client.NewTimeFrame(d)
	tf, err := client.UpdateTimeFrame(ctx, c, id, body)
	if err != nil {
		return diag.FromErr(err)
	}
	return timeFrameToResource(d, tf)
}

func timeFrameDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) (diags diag.Diagnostics) {
	c := meta.(*client.Client)
	id := d.Id()
	_, err := client.DeleteTimeFrame(ctx, c, id)
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
