---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pfptmeta_group Resource - terraform-provider-pfptmeta"
subcategory: ""
description: |-
  Groups represent a collection of users, typically belong to a common department or share same privileges in the organization.
---

# pfptmeta_group (Resource)

Groups represent a collection of users, typically belong to a common department or share same privileges in the organization.

## Example Usage

```terraform
resource "pfptmeta_group" "new_group" {
  name        = "group name"
  description = "group description"
  expression  = "tag_name:tag_value OR platform:macOS"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String)

### Optional

- **description** (String)
- **expression** (String) Allows grouping entities by their tags. Filtering by tag value is also supported if provided. Supported operations: AND, OR, parenthesis.

### Read-Only

- **id** (String) The ID of this resource.

