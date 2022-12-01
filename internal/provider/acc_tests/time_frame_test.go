package acc_tests

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"regexp"
	"testing"
)

const (
	timeFrameStep1 = `
resource "pfptmeta_time_frame" "tf" {
  name        = "tf"
  description = "tf desc"
  days        = ["monday"]
  start_time {
    hour   = 8
    minute = 0
  }
  end_time {
    hour   = 18
    minute = 0
  }
}
`
	timeFrameStep2 = `
resource "pfptmeta_time_frame" "tf" {
  name        = "tf 1"
  description = "tf desc 1"
  days        = ["tuesday"]
  start_time {
    hour   = 9
    minute = 30
  }
  end_time {
    hour   = 17
    minute = 30
  }
}
`
	timeFrameStepDataSource = `
data "pfptmeta_time_frame" "in" {
  id = pfptmeta_time_frame.tf.id
}
`
)

func TestAccResourceTimeFrame(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      validateResourceDestroyed("time_frame", "v1/time_frames"),
		Steps: []resource.TestStep{
			{
				Config: timeFrameStep1,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_time_frame.tf", "id", regexp.MustCompile("^tmf-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "name", "tf"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "description", "tf desc"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "days.0", "monday"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "start_time.0.hour", "8"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "start_time.0.minute", "0"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "end_time.0.hour", "18"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "end_time.0.minute", "0"),
				),
			},
			{
				Config: timeFrameStep2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_time_frame.tf", "id", regexp.MustCompile("^tmf-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "name", "tf 1"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "description", "tf desc 1"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "days.0", "tuesday"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "start_time.0.hour", "9"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "start_time.0.minute", "30"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "end_time.0.hour", "17"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "end_time.0.minute", "30"),
				),
			},
		},
	})
}

func TestAccDataSourceTimeFrame(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: timeFrameStep1 + timeFrameStepDataSource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("pfptmeta_time_frame.tf", "id", regexp.MustCompile("^tmf-.+$")),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "name", "tf"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "description", "tf desc"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "days.0", "monday"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "start_time.0.hour", "8"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "start_time.0.minute", "0"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "end_time.0.hour", "18"),
					resource.TestCheckResourceAttr("pfptmeta_time_frame.tf", "end_time.0.minute", "0"),
				),
			},
		},
	})
}
