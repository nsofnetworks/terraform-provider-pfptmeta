resource "pfptmeta_enterprise_dns" "enterprise_dns" {
  name        = "enterprise dns name"
  description = "enterprise dns description"
  mapped_domains {
    name          = "mapped.domain1.com"
    mapped_domain = "mapped.domain1.com"
  }
  mapped_domains {
    name          = "mapped.domain2.com"
    mapped_domain = "mapped.domain2.com"
  }
}