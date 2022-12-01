resource "pfptmeta_time_frame" "tf" {
  name        = "time frame name"
  description = "time frame description"
  days        = ["monday", "tuesday", "wednesday", "thursday", "friday"]
  start_time {
    hour   = 8
    minute = 0
  }
  end_time {
    hour   = 18
    minute = 0
  }
}