---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Resource pfptmeta_device - terraform-provider-pfptmeta"
subcategory: "Network Resources"
description: |-
  When a user is onboarded to the Proofpoint NaaS platform via Proofpoint Agent,
  the user identity is bound to the device of the logging request, and a certificate is issued to this machine.
---

# Resource (pfptmeta_device)

When a user is onboarded to the Proofpoint NaaS platform via Proofpoint Agent,
the user identity is bound to the device of the logging request, and a certificate is issued to this machine.

## Example Usage

```terraform
resource "pfptmeta_device" "device" {
  name        = "my device"
  description = "some details about the device"
  owner_id    = "usr-abc123"
  tags = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)

### Optional

- `description` (String)
- `enabled` (Boolean)
- `owner_id` (String)
- `tags` (Map of String) Key/value attributes for combining elements together into Smart Groups, and placed as targets or sources in Policies

### Read-Only

- `auto_aliases` (List of String)
- `groups` (List of String)
- `id` (String) The ID of this resource.
