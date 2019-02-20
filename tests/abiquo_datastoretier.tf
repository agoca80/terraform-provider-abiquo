resource "abiquo_datastoretier" "test" {
  datacenter  = "${data.abiquo_datacenter.test.id}"
  description = "required description"
  enabled     = true
  name        = "test dstier"
  policy      = "PERFORMANCE"
}

data "abiquo_datacenter" "test" { name = "${var.datacenter}" }
