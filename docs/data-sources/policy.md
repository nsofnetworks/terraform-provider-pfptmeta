---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pfptmeta_policy Data Source - terraform-provider-pfptmeta"
subcategory: ""
description: |-
  Policies bind network elements (devices, services and subnets) to users, groups and to other network elements and define access direction and connections.
---

# pfptmeta_policy (Data Source)

Policies bind network elements (devices, services and subnets) to users, groups and to other network elements and define access direction and connections.

## Example Usage

```terraform
data "pfptmeta_policy" "policy" {
  id = "pol-acb123"
}

output "policy" {
  value = data.pfptmeta_policy.policy
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **id** (String) The ID of this resource.

### Read-Only

- **description** (String)
- **destinations** (Set of String) Entities (users, groups or network elements) to which the access is granted to.
- **enabled** (Boolean)
- **exempt_sources** (Set of String) Entities (users, groups or network elements) to be excluded from accessing the application defined in this policy.
- **name** (String)
- **protocol_groups** (Set of String) Protocol groups that restrict the protocols or TCP/UDP ports for this policy
- **sources** (Set of String) Entities (users, groups or network elements) to be authorized to access the application defined in this policy.

