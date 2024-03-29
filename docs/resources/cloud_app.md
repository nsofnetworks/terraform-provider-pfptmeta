---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Resource pfptmeta_cloud_app - terraform-provider-pfptmeta"
subcategory: "Web Security Resources"
description: |-
  Additional flexibility can be achieved by controlling user access to cloud-based applications.
  The administrators can block, redirect to isolation or record any attempt to access applications selected from a
  vast Proofpoint catalog or defined via a specific URL.
---

# Resource (pfptmeta_cloud_app)

Additional flexibility can be achieved by controlling user access to cloud-based applications.
The administrators can block, redirect to isolation or record any attempt to access applications selected from a 
vast Proofpoint catalog or defined via a specific URL.

## Example Usage

```terraform
resource "pfptmeta_cloud_app" "ca" {
  name        = "cloud app"
  description = "cloud app description"
  app         = "sia-abc123"
  urls        = ["192.6.6.5", "ynet.co.il"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)

### Optional

- `app` (String) The ID of the [catalog_app](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/data-sources/catalog_app) data-source.
- `description` (String)
- `tenant` (String) Specific tenant ID of the app on which the cloud application rule should be applied. 
Valid only for catalog apps that have [tenant_corp_id_support](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/data-sources/catalog_app#tenant_corp_id_support) set to true
- `tenant_type` (String) ENUM: `All`, `Personal`, `Corporate` (Defaults to All). Valid only for catalog apps that have [tenant_type_support](https://registry.terraform.io/providers/nsofnetworks/pfptmeta/latest/docs/data-sources/catalog_app#tenant_type_support) set to true
- `urls` (List of String) A list of URLs to associate with this cloud app.

### Read-Only

- `id` (String) The ID of this resource.
- `type` (String)
