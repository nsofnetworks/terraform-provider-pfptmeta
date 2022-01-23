resource "pfptmeta_trusted_network" "network" {
  name         = "trusted network name"
  description  = "trusted network description"
  apply_to_org = true
  criteria {
    external_ip_config {
      addresses_ranges = ["192.1.0.0/16"]
    }
  }
  criteria {
    resolved_address_config {
      addresses_ranges = ["192.1.0.0/16"]
      hostname         = "office.address11.com"
    }
  }
}