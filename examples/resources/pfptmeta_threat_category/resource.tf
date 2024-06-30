resource "pfptmeta_threat_category" "tc" {
  name             = "threat category"
  description      = "threat category description"
  confidence_level = "HIGH"
  risk_level       = "MEDIUM"
  countries        = ["AF", "EG"]
  types            = ["Peer to Peer", "Scanner"]
  third_party_app  = ["MALICIOUS"]
}

