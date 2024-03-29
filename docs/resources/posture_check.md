---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "Resource pfptmeta_posture_check - terraform-provider-pfptmeta"
subcategory: "Device Management"
description: |-
  Posture checks are administrator-defined sets of criteria allowing to or preventing the user devices from connecting to Proofpoint NaaS.
  Administrators can use SQL to create real-world conditions based on underlying information from the operating system and its hardware or choose from a list of common pre-defined conditions.
  For example, the presence (or absence) of a specified file or a system process can serve as a pre-condition for letting a device to access the Proofpoint NaaS.
  The posture checks can be based on SQL query strings. The queries use the osquery framework, see osquery.io https://osquery.io/ for details on osquery.
  Posture checks can be viewed and filtered by failure via security logs see here https://help.metanetworks.com/knowledgebase/posture_checks for more details.
---

# Resource (pfptmeta_posture_check)

Posture checks are administrator-defined sets of criteria allowing to or preventing the user devices from connecting to Proofpoint NaaS.
Administrators can use SQL to create real-world conditions based on underlying information from the operating system and its hardware or choose from a list of common pre-defined conditions.
For example, the presence (or absence) of a specified file or a system process can serve as a pre-condition for letting a device to access the Proofpoint NaaS.
The posture checks can be based on SQL query strings. The queries use the osquery framework, see [osquery.io](https://osquery.io/) for details on osquery.
Posture checks can be viewed and filtered by failure via security logs see [here](https://help.metanetworks.com/knowledgebase/posture_checks) for more details.

## Example Usage

```terraform
resource "pfptmeta_user" "user" {
  given_name  = "John"
  family_name = "Smith"
  email       = "john.smith@example.com"
}

resource "pfptmeta_network_element" "device" {
  name     = "device-name"
  owner_id = pfptmeta_user.user.id
  platform = "Linux"
}

resource "pfptmeta_posture_check" "antivirus_check" {
  name                 = "CrowdStrike Posture Check1"
  description          = "Description"
  apply_to_entities    = [pfptmeta_network_element.device.id, pfptmeta_user.user.id]
  osquery              = "select * from processes where name='falcon-sensor' and state='S';"
  platform             = "Linux"
  enabled              = true
  action               = "NONE"
  when                 = ["PERIODIC", "PRE_CONNECT"]
  interval             = 60
  user_message_on_fail = "CrowdStrike is not installed on device1"
}

resource "pfptmeta_posture_check" "min_client_version" {
  name            = "Min Client Version1"
  apply_to_org    = true
  exempt_entities = [pfptmeta_user.user.id]
  check {
    type        = "minimum_app_version"
    min_version = "4.0.0"
  }
  platform             = "iOS"
  action               = "DISCONNECT"
  when                 = ["PRE_CONNECT"]
  user_message_on_fail = "CrowdStrike is not installed on device1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String)
- `when` (List of String) When to run the check, ENUM: `PRE_CONNECT`, `PERIODIC`.

### Optional

- `action` (String) Action to take in case a posture check fails. ENUM: `DISCONNECT`, `NONE`, `WARNING`:
	- **Disconnect** - disconnect device from Proofpoint NaaS.
	- **None** - do nothing, useful during the discovery phase. 
   - **Warning** - pop up a warning message, useful during the discovery phase.
- `apply_to_entities` (List of String) Entities (users, groups or devices) to be applied in the posture check.
- `apply_to_org` (Boolean) Whether to apply to all devices on the organization. Note: this attribute overrides `apply_to_entities`
- `check` (Block List, Max: 1) Predefined checks. cannot be set with `osquery`. (see [below for nested schema](#nestedblock--check))
- `description` (String)
- `enabled` (Boolean) Defaults to true
- `exempt_entities` (List of String) Entities (users, groups or devices) which are exempt from the posture check.
- `interval` (Number) Interval in minutes between checks, mandatory when `when` is set to `PERIODIC`. ENUM: 5, 60.
- `osquery` (String) osquery to use in the posture check, see [here](https://osquery.io/) for more details.
- `platform` (String) Device platforms that should be applied in the posture check. ENUM: `Android`, `macOS`, `iOS`, `Linux`, `Windows`, `ChromeOS`.
- `user_message_on_fail` (String) Message to be displayed when posture check fails.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--check"></a>
### Nested Schema for `check`

Required:

- `type` (String) ENUM: `jailbroken_rooted`, `screen_lock_enabled`, `minimum_app_version`, `minimum_os_version`, `malicious_app_detection`, `developer_mode_enabled`.

Optional:

- `min_version` (String) Minimum version required by the check. Required when `type` is `minimum_app_version` or `minimum_os_version`, format: major.minor.patch.
