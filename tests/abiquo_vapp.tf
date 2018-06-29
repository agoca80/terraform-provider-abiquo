data "abiquo_vdc"        "test"       { name = "tests" }

resource "abiquo_vapp" "test" {
  virtualdatacenter = "${data.abiquo_vdc.test.id}"
  name              = "testAccVappBasic"
}
