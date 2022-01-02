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