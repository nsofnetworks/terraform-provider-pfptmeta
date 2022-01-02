---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "pfptmeta_certificate Resource - terraform-provider-pfptmeta"
subcategory: "Administration"
description: |-
  SSL certificate. It is used mostly to allow EasyLinks to utilize HTTPS, when operating with the redirect or native access types.
---

# pfptmeta_certificate (Resource)

SSL certificate. It is used mostly to allow EasyLinks to utilize HTTPS, when operating with the `redirect` or `native` access types.

## Example Usage

```terraform
resource "pfptmeta_certificate" "cert" {
  name        = "certificate name"
  description = "certificate description"
  sans        = ["test.example.com"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **name** (String)
- **sans** (Set of String) List of certificate SANs

### Optional

- **description** (String)

### Read-Only

- **id** (String) The ID of this resource.
- **serial_number** (String)
- **status** (String) Certificate state, can be one of the following:
	- **Pending** - Initial state that may take several minutes. During this stage, a request is sent to certification authority and the system is waiting for the certificate approval.
	- **OK** - Certificate has been validated by the certification authority and ready for use.
	- **Warning** - Certificate is valid, but it is to expire within 30 days. DNS check attempts for the certificate renewal have failed.
	- **Error** - Certificate has expired, all DNS checks have failed so far, and no renewal attempts are being made.
- **status_description** (String)
- **valid_not_after** (String)
- **valid_not_before** (String)