# resource "pfptmeta_device" "device" {
#   name        = "my device"
#   description = "some details about the device"
#   owner_id    = "usr-abc123"
#   tags = {
#     tag_name1 = "tag_value1"
#     tag_name2 = "tag_value2"
#   }
# }


terraform {
  required_providers {
    pfptmeta = {
      source  = "nsofnetworks/pfptmeta"
      version = "0.1.49"
    }
  }
}

provider "pfptmeta" {
  api_key       = "key-xe35s61AkYyDMrj"
  api_secret    = "cb8f94a3ca32fe8a1d2eb254b44d2a2c870b9588845a5206be18aa465f89e64ef6996c1e7c5347ee"
  org_shortname = "playground"
  realm         = "us"
}

resource "pfptmeta_network_element" "mapped_service2" {
  name           = "example terraform"
  description    = "some details about the mapped service"
  mapped_service = "mapped.service.com"
  tags = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
}
