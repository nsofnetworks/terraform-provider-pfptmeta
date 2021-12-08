data "pfptmeta_policy" "policy" {
  id = "pol-acb123"
}

output "policy" {
  value = data.pfptmeta_policy.policy
}