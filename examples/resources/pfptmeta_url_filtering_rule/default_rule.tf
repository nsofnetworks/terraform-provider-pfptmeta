resource "pfptmeta_threat_category" "malicious" {
  name             = "Malicious Threat"
  confidence_level = "LOW"
  risk_level       = "LOW"
  countries        = ["IR", "KP"]
  types = [
    "Bitcoin Related", "Blackhole", "Botnets", "Brute Forcer", "CnC", "Compromised", "Drop", "EXE Source",
    "Fake AV", "Keyloggers and Monitoring", "Malware Sites", "Mobile CnC", "Mobile Spyware CnC", "P2P CnC",
    "Phishing and Other Frauds", "Spyware and Adware", "Tor"
  ]
}

resource "pfptmeta_content_category" "strict" {
  name             = "Strict Category"
  confidence_level = "LOW"
  types = [
    "Sex Education", "Nudity", "Abused Drugs", "Marijuana", "Swimsuits and Intimate Apparel", "Violence",
    "Gross", "Adult and Pornography", "Weapons", "Hate and Racism", "Gambling"
  ]
  urls = [".espn.com"]
}

resource "pfptmeta_url_filtering_rule" "default_rule" {
  name                         = "default rule"
  description                  = "default rule"
  apply_to_org                 = true
  action                       = "BLOCK"
  advanced_threat_protection   = true
  threat_categories            = [pfptmeta_threat_category.malicious.id]
  forbidden_content_categories = [pfptmeta_content_category.strict.id]
  priority                     = 94
  warn_ttl                     = 15
}