---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Resource pfptmeta_proxy_port_range - terraform-provider-pfptmeta"
subcategory: "Web Security Resources"
description: |-
  Administrators can define communication ports for Web Security traffic over HTTP/S. The following port range is supported: 1 – 65535.
---

# Resource (pfptmeta_proxy_port_range)

Administrators can define communication ports for Web Security traffic over HTTP/S. The following port range is supported: 1 – 65535.

## Example Usage

```terraform
resource "pfptmeta_proxy_port_range" "proxy_port_range" {
  name        = "my port range"
  description = "some details about destination port ranges"
  proto       = "HTTP"
  from_port   = 20000
  to_port     = 20100
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `from_port` (Number) Start port for this proxy port range
- `name` (String)
- `proto` (String) Protocol for this proxy port range. ENUM: `HTTPS`,`HTTP`
- `to_port` (Number) End port for this proxy port range

### Optional

- `description` (String)
- `read_only` (Boolean)

### Read-Only

- `id` (String) The ID of this resource.
