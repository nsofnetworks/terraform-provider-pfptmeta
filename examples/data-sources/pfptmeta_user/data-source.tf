data "pfptmeta_user" "user_by_id" {
  id = "usr-abc123"
}

data "pfptmeta_user" "user_by_email" {
  email = "John.Smith@example.com"
}

output "user_by_id" {
  value = data.pfptmeta_user.user_by_id
}

output "user_by_email" {
  value = data.pfptmeta_user.user_by_email
}