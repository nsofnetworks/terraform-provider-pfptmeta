# data "pfptmeta_trusted_network" "network" {
#   id = "tn-123abc"
# }

# output "network" {
#   value = data.pfptmeta_trusted_network.network
# }

# terraform {
#   required_providers {
#     pfptmeta = {
#       source  = "nsofnetworks/pfptmeta"
#       version = "0.1.39"
#     }
#   }
# }

# provider "pfptmeta" {
#   api_key       = "key-xe35s61AkYyDMrj"
#   api_secret    = "2ea1557b6cca42e6052966a6016b2f074fe8daf08ac1205ba5fce6533fd67370a6e3752f672a1dc1"
#   org_shortname = "playground"
#   realm         = "us"
# }

# provider "pfptmeta" {
#   api_key       = "key-KjUnoNz81ZkLo"
#   api_secret    = "bac628ebbe43689d9fa3a0b121c5f720f1b85dbf581e320c8adced8562e879a1a5e6bef364ccde2b"
#   org_shortname = "nsof"
# }


resource "pfptmeta_trusted_network" "gadi-network" {
  name         = "gadi-test-tn"
  description  = "trusted network description"
  apply_to_org = false
  apply_to_entities = ["ne-656"]
}