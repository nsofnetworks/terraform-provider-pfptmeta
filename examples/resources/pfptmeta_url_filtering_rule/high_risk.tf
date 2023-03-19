data "pfptmeta_catalog_app" "salesforce" {
  name     = "Salesforce"
  category = "Content Management"
}

resource "pfptmeta_cloud_app" "salesforce" {
  name = "salesforce"
  app  = data.pfptmeta_catalog_app.salesforce.id
  urls = ["my.salesforce.com"]
}

resource "pfptmeta_url_filtering_rule" "high_risk" {
  name              = "Block High Risk - Expression"
  apply_to_org      = true
  action            = "BLOCK"
  cloud_apps        = [pfptmeta_cloud_app.salesforce.id]
  filter_expression = "crwd_agent:fail OR crwdzta:high OR npre-it"
  priority          = 90
  warn_ttl          = 15
}
