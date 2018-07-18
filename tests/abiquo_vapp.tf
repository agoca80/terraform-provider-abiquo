resource "abiquo_vapp" "test" {
  virtualdatacenter = "${data.abiquo_vdc.test.id}"
  name              = "test vapp"
}

data "abiquo_vdc"  "test" { name = "tests" }
data "abiquo_vapp" "test" {
  virtualdatacenter = "${data.abiquo_vdc.test.id}"
  name              = "${abiquo_vapp.test.name}"
}
