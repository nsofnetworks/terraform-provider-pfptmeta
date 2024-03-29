---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Resource pfptmeta_tenant_restriction - terraform-provider-pfptmeta"
subcategory: "Web Security Resources"
description: |-
  Tenant restrictions are used to control access to SaaS cloud applications, such as Office 365. This is done to prevent unwarranted access to potentially malicious tenants, or the ones that are forbidden from accessing due to risk of data loss. With tenant restrictions, organizations can specify the list of tenants that their users are permitted to access, blocking all other tenants. When tenant restrictions are used, the Web Security proxy inserts a list of permitted tenants into traffic destined for the relevant cloud app, allowing it to perform the enforcement.
---

# Resource (pfptmeta_tenant_restriction)

Tenant restrictions are used to control access to SaaS cloud applications, such as Office 365. This is done to prevent unwarranted access to potentially malicious tenants, or the ones that are forbidden from accessing due to risk of data loss. With tenant restrictions, organizations can specify the list of tenants that their users are permitted to access, blocking all other tenants. When tenant restrictions are used, the Web Security proxy inserts a list of permitted tenants into traffic destined for the relevant cloud app, allowing it to perform the enforcement.

## Example Usage

```terraform
resource "pfptmeta_tenant_restriction" "google" {
  name        = "google"
  description = "google tenant restriction"
  google_config {
    allow_consumer_access  = true
    allow_service_accounts = false
    tenants                = ["altostrat.com", "tenorstrat.com"]
  }
}

resource "pfptmeta_tenant_restriction" "microsoft" {
  name        = "microsoft"
  description = "microsoft tenant restriction"
  microsoft_config {
    allow_personal_microsoft_domains = true
    tenant_directory_id              = "456ff232-35l2-5h23-b3b3-3236w0826f3d"
    tenants                          = ["onmicrosoft.com"]
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)

### Optional

- `description` (String)
- `google_config` (Block List, Max: 1) (see [below for nested schema](#nestedblock--google_config))
- `microsoft_config` (Block List, Max: 1) (see [below for nested schema](#nestedblock--microsoft_config))

### Read-Only

- `id` (String) The ID of this resource.
- `type` (String)

<a id="nestedblock--google_config"></a>
### Nested Schema for `google_config`

Required:

- `allow_consumer_access` (Boolean) Whether to allow access consumer Google Accounts, such as @gmail.com and @googlemail.com.
- `allow_service_accounts` (Boolean) Whether to allow access to authenticated service accounts.
- `tenants` (List of String) Configuring this will cause Google to issue security token only for the specified tenants.


<a id="nestedblock--microsoft_config"></a>
### Nested Schema for `microsoft_config`

Required:

- `allow_personal_microsoft_domains` (Boolean) Whether to allow Microsoft applications for consumer accounts.
- `tenant_directory_id` (String) The directory ID of the tenant that sets tenant restrictions.
- `tenants` (List of String) Configuring this will cause Azure AD to issue security token only for the specified tenants. Any domain that is registered with a tenant can be used to identify the tenant in this list, as well as the directory ID itself
