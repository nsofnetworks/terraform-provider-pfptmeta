data "pfptmeta_location" "new_york" {
  name = "LGA"
}

output "new_york_location" {
  value = data.pfptmeta_location.new_york
}
