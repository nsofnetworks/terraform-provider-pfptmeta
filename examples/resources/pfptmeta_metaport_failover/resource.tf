resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "ms"
  mapped_subnets = ["0.0.0.0/0"]
}

resource "pfptmeta_metaport" "metaport1" {
  name = "metaport1"
}

resource "pfptmeta_metaport" "metaport2" {
  name = "metaport2"
}

resource "pfptmeta_metaport_cluster" "metaport_cluster1" {
  name      = "metaport cluster1"
  metaports = [pfptmeta_metaport.metaport1.id]
}

resource "pfptmeta_metaport_cluster" "metaport_cluster2" {
  name      = "metaport cluster2"
  metaports = [pfptmeta_metaport.metaport2.id]
}

resource "pfptmeta_notification_channel" "mail" {
  name        = "mail-channel"
  description = "mail channel description"
  email_config {
    recipients = ["user1@example.com", "user2@example.com"]
  }
}

resource "pfptmeta_metaport_failover" "failover" {
  name                  = "metaport failover name"
  description           = "metaport failover description"
  mapped_elements       = [pfptmeta_network_element.mapped-subnet.id]
  cluster_1             = pfptmeta_metaport_cluster.metaport_cluster1
  cluster_2             = pfptmeta_metaport_cluster.metaport_cluster2
  notification_channels = [pfptmeta_notification_channel.mail.id]
  failback {
    trigger = "auto"
  }
  failover {
    delay     = 1
    threshold = 1
    trigger   = "auto"
  }
}