data "pfptmeta_time_frame" "tf" {
  id = "tf-123abc"
}

output "time_frame" {
  value = data.pfptmeta_time_frame.tf
}