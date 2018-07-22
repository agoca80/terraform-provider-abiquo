resource "abiquo_datastoreloadrule" "test" {
  load   = "95"
  target = "${data.abiquo_datacenter.test.id}"
}

data "abiquo_datacenter" "test" { name = "datacenter 1" }
