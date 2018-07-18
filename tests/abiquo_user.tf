resource "abiquo_user" "test" {
  enterprise = "${data.abiquo_enterprise.test.id}"
  role       = "${data.abiquo_role.test.id}"
  active     = true
  name       = "test"
  surname    = "test"
  nick       = "test"
  email      = "test@test.com"
}

data "abiquo_enterprise" "test" { name = "Abiquo" }
data "abiquo_role"       "test" { name = "CLOUD_ADMIN" }
