resource "pfptmeta_group" "new_group" {
  name = "easylink-group"
}

locals {
  hostname = "test.example.com"
  ipv4     = "196.10.10.1"
}

resource "pfptmeta_certificate" "cert" {
  name = "certificate name"
  sans = [local.hostname]
}

resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service name"
  mapped_service = local.ipv4
}

resource "pfptmeta_network_element_alias" "alias" {
  network_element_id = pfptmeta_network_element.mapped-service.id
  alias              = local.hostname
}

resource "pfptmeta_easylink" "meta_easylink" {
  name        = "meta easylink name"
  description = "meta easylink description"
  domain_name = local.hostname
  access_type = "meta"
  port        = 443
  protocol    = "https"
  viewers     = [pfptmeta_group.new_group.id]
}

resource "pfptmeta_easylink" "meta_rdp_easylink" {
  name        = "meta_rdp easylink name"
  description = "meta_rdp easylink description"
  domain_name = local.hostname
  access_type = "meta"
  port        = 3389
  protocol    = "rdp"
  viewers     = [pfptmeta_group.new_group.id]
  rdp {
    security               = "nla"
    server_keyboard_layout = "french"
  }
}

resource "pfptmeta_easylink" "redirect_easylink" {
  name              = "redirect easylink name"
  description       = "redirect easylink description"
  domain_name       = local.ipv4
  access_fqdn       = local.hostname
  access_type       = "redirect"
  port              = 443
  protocol          = "https"
  mapped_element_id = pfptmeta_network_element.mapped-service.id
  viewers           = [pfptmeta_group.new_group.id]
  certificate_id    = pfptmeta_certificate.cert.id
  root_path         = "/application"
}

resource "pfptmeta_easylink" "native_easylink" {
  name              = "native easylink name"
  description       = "native easylink description"
  domain_name       = local.ipv4
  access_fqdn       = local.hostname
  access_type       = "native"
  port              = 443
  protocol          = "https"
  mapped_element_id = pfptmeta_network_element.mapped-service.id
  viewers           = [pfptmeta_group.new_group.id]
  certificate_id    = pfptmeta_certificate.cert.id
  root_path         = "/application"
  proxy {
    rewrite_content_types = ["json", "html"]
    rewrite_http          = true
  }
}



