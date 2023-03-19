resource "pfptmeta_cloud_app" "ca" {
  name        = "cloud app"
  description = "cloud app description"
  app         = "sia-abc123"
  urls        = ["192.6.6.5", "ynet.co.il"]
}
