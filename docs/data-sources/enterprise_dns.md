---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pfptmeta_enterprise_dns Data Source - terraform-provider-pfptmeta"
subcategory: "Network Resources"
description: |-
  Enterprise DNS provides integration with global, enterprise DNS servers, allowing resolution of FQDNs for domains that are in different locations/datacenters.
---

# pfptmeta_enterprise_dns (Data Source)

Enterprise DNS provides integration with global, enterprise DNS servers, allowing resolution of FQDNs for domains that are in different locations/datacenters.

## Example Usage

```terraform
data "pfptmeta_enterprise_dns" "enterprise_dns" {
  id = "ed-123"
}

output "enterprise_dns" {
  value = data.pfptmeta_enterprise_dns.enterprise_dns
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **id** (String) The ID of this resource.

### Read-Only

- **description** (String)
- **mapped_domains** (Set of Object) DNS suffixes to be resolved within the enterprise DNS server (see [below for nested schema](#nestedatt--mapped_domains))
- **name** (String)

<a id="nestedatt--mapped_domains"></a>
### Nested Schema for `mapped_domains`

Read-Only:

- **mapped_domain** (String)
- **name** (String)
