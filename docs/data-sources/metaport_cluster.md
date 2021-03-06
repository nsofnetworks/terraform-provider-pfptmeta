---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pfptmeta_metaport_cluster Data Source - terraform-provider-pfptmeta"
subcategory: "Network"
description: |-
  MetaPort cluster defines a group of highly-available MetaPorts that are deployed together in a single data center
---

# pfptmeta_metaport_cluster (Data Source)

MetaPort cluster defines a group of highly-available MetaPorts that are deployed together in a single data center

## Example Usage

```terraform
data "pfptmeta_metaport_cluster" "metaport_cluster_by_id" {
  id = "mpc-123"
}

output "metaport_cluster_by_id" {
  value = data.pfptmeta_metaport_cluster.metaport_cluster_by_id
}

data "pfptmeta_metaport_cluster" "metaport_cluster_by_name" {
  name = "metaport cluster name"
}

output "metaport_cluster_by_name" {
  value = data.pfptmeta_metaport_cluster.metaport_cluster_by_name
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **id** (String) The ID of this resource.
- **name** (String)

### Read-Only

- **description** (String)
- **mapped_elements** (Set of String) List of mapped element IDs
- **metaports** (Set of String) List of MetaPort IDs
