resource "abiquo_virtualappliance" "test" {
  virtualdatacenter = "${data.abiquo_virtualdatacenter.test.id}"
  name              = "test vapp"
}

data "abiquo_virtualdatacenter"  "test" { name = "tests" }
data "abiquo_virtualappliance" "test" {
  virtualdatacenter = "${data.abiquo_virtualdatacenter.test.id}"
  name              = "${abiquo_virtualappliance.test.name}"
}
