resource "pfptmeta_tenant_restriction" "google" {
  name        = "google"
  description = "google tenant restriction"
  google_config {
    allow_consumer_access  = true
    allow_service_accounts = false
    tenants                = ["altostrat.com", "tenorstrat.com"]
  }
}

resource "pfptmeta_tenant_restriction" "microsoft" {
  name        = "microsoft"
  description = "microsoft tenant restriction"
  microsoft_config {
    allow_personal_microsoft_domains = true
    tenant_directory_id              = "456ff232-35l2-5h23-b3b3-3236w0826f3d"
    tenants                          = ["onmicrosoft.com"]
  }
}