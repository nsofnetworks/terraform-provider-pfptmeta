resource "pfptmeta_threat_category" "log_all" {
  name             = "Threats To Log"
  confidence_level = "LOW"
  risk_level       = "LOW"
  types = [
    "Botnets", "Chat Server", "Phishing and Other Frauds", "Utility", "Self Signed SSL", "Brute Forcer", "DDoS Target",
    "Parking", "Scanner", "Online Gaming", "Blackhole", "Compromised", "Fake AV", "Malware Sites", "P2P CnC",
    "Remote Access Service", "Bitcoin Related", "SPAM URLs", "Mobile CnC", "Tor", "IP Check", "DynDNS", "CnC",
    "Spyware and Adware", "Undesirable", "Mobile Spyware CnC", "Abused TLD", "EXE Source", "VPN", "Drop",
    "Proxy Avoidance and Anonymizers", "Peer to Peer"
  ]
}

resource "pfptmeta_content_category" "log_all" {
  name             = "Category To Log"
  confidence_level = "HIGH"
  types = [
    "Abortion", "Sex Education", "Pay to Surf", "Web Advertisements", "Dynamically Generated Content",
    "Parked Domains", "Alcohol and Tobacco", "Personal sites and Blogs", "Hacking",
    "Abused Drugs", "Marijuana", "Training and Tools", "Reference and Research", "Educational Institutions",
    "Web-based Email", "Financial Services", "Business and Economy",
    "Individual Stock Advice and Tools", "Home and Garden", "Gambling", "Games", "Kids", "Legal", "Government",
    "Health and Medicine", "Recreation and Hobbies", "Questionable", "Cheating", "Illegal", "Job Search",
    "Swimsuits and Intimate Apparel", "Hate and Racism", "Local Information", "News and Media", "Nudity",
    "Philosophy and Political Advocacy", "Adult and Pornography", "Internet Portals", "Real Estate", "Cult and Occult",
    "Religion", "Search Engines", "Image and Video Search", "Auctions", "Shopping", "Online Greeting Cards",
    "Fashion and Beauty", "Social Networking", "Dating", "Society", "Computer and Internet Security",
    "Computer and Internet Info", "Shareware and Freeware", "Personal Storage", "Content Delivery Networks",
    "Web Hosting", "Internet Communications", "Hunting and Fishing", "Sports", "Streaming Media",
    "Entertainment and Arts", "Translation", "Travel", "Motor Vehicles", "Violence", "Gross", "Weapons"
  ]
  urls = ["clarivate.io"]
}

resource "pfptmeta_url_filtering_rule" "log_all" {
  name                         = "Log All"
  apply_to_org                 = true
  action                       = "LOG"
  advanced_threat_protection   = true
  threat_categories            = [pfptmeta_threat_category.log_all.id]
  forbidden_content_categories = [pfptmeta_content_category.log_all.id]
  priority                     = 95
  warn_ttl                     = 15
}
