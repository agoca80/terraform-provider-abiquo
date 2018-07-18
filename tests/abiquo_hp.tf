resource "abiquo_hp" "test" {
  active = true
  name = "test hp"
  cpu  = 16
  ram  = 64
  datacenter = "${data.abiquo_datacenter.test.id}"
}

data "abiquo_datacenter" "test" { name = "datacenter 1" }
