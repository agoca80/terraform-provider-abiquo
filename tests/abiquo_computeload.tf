resource "abiquo_computeload" "test" {
  cpuload    = "1000"
  ramload    = "95"
  target     = "${data.abiquo_datacenter.test.id}"
}

data "abiquo_datacenter" "test" { name = "datacenter 1" }
