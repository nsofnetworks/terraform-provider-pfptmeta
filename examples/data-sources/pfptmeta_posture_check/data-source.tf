data "pfptmeta_posture_check" "check" {
  id = "pc-123abc"
}

output "check" {
  value = data.pfptmeta_posture_check.check
}