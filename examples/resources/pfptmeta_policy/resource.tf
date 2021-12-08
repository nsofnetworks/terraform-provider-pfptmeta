data "pfptmeta_protocol_group" "HTTPS" {
  name = "HTTPS"
}

data "pfptmeta_group" "developers_group" {
  name = "dev"
}

resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service"
  mapped_service = "mapped.service.com"
}

resource "pfptmeta_policy" "policy" {
  name            = "policy name"
  description     = "policy description"
  sources         = [data.pfptmeta_group.developers_group.id]
  destinations    = [pfptmeta_network_element.mapped-service.id]
  protocol_groups = [data.pfptmeta_protocol_group.HTTPS.id]
}