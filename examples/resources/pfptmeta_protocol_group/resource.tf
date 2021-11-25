resource "pfptmeta_protocol_group" "new_protocol" {
  name = "NEW_PROTOCOL"
  protocols {
    from_port = 445
    to_port   = 445
    proto     = "udp"
  }
  protocols {
    from_port = 446
    to_port   = 446
    proto     = "tcp"
  }
}