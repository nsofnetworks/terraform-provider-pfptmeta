resource "pfptmeta_user" "user" {
  given_name  = "John"
  family_name = "Smith"
  email       = "john.smith@example.com"
  description = "some details about the user"
  phone       = "+97251234567"
  tags = {
    tag_name1 = "tag_value1"
    tag_name2 = "tag_value2"
  }
}