resource "pfptmeta_group" "group" {
  name = "some-group"
}

resource "pfptmeta_user" "user1" {
  given_name  = "user"
  family_name = "one"
  email       = "user1@example.com"
}

resource "pfptmeta_user" "user2" {
  given_name  = "user"
  family_name = "two"
  email       = "user2@example.com"
}

resource "pfptmeta_user" "user3" {
  given_name  = "user"
  family_name = "three"
  email       = "user3@example.com"
}

resource "pfptmeta_group_users_attachment" "attachment" {
  group_id = pfptmeta_group.group.id
  users = [
    pfptmeta_user.user1.id,
    pfptmeta_user.user2.id,
    pfptmeta_user.user3.id
  ]
}