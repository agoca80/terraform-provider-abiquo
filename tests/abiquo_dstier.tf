data "abiquo_datacenter" "test" { name = "datacenter 1" }

resource "abiquo_dstier" "test" {
  datacenter  = "${data.abiquo_datacenter.test.id}"
  description = "required description"
  enabled     = true
  name        = "testAccDSTierBasic"
  policy      = "PERFORMANCE"
}
