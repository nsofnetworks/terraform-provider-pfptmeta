resource "pfptmeta_content_category" "news" {
  name             = "News"
  confidence_level = "LOW"
  types            = ["News and Media"]
}

resource "pfptmeta_content_category" "social_network" {
  name             = "Social Networking"
  confidence_level = "LOW"
  types            = ["Social Networking"]
}

resource "pfptmeta_time_frame" "work_hours" {
  name = "Work Hours"
  days = ["monday", "tuesday", "wednesday", "thursday", "friday"]
  start_time {
    hour   = 8
    minute = 0
  }
  end_time {
    hour   = 18
    minute = 0
  }
}

resource "pfptmeta_url_filtering_rule" "work_time" {
  name                         = "News And Social Networking"
  description                  = "Blocks news and social networking during work hours"
  apply_to_org                 = true
  action                       = "BLOCK"
  advanced_threat_protection   = false
  forbidden_content_categories = [pfptmeta_content_category.news.id, pfptmeta_content_category.social_network.id]
  priority                     = 80
  warn_ttl                     = 15
  schedule                     = [pfptmeta_time_frame.work_hours.id]
}