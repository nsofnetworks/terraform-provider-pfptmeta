resource "pfptmeta_network_element" "mapped-subnet" {
  name           = "mapped subnet name"
  description    = "some details about the mapped subnet"
  mapped_subnets = ["0.0.0.0/0", "10.20.30.0/24"]
  tags = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
  mapped_domains {
    name          = "test.com"
    mapped_domain = "test.com"
  }
  mapped_domains {
    name          = "test1.com"
    mapped_domain = "test1.com"
  }
}

resource "pfptmeta_network_element" "mapped-service" {
  name           = "mapped service name"
  description    = "some details about the mapped service"
  mapped_service = "mapped.service.com"
  tags = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
}