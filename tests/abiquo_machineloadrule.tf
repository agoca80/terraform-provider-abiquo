resource "abiquo_machineloadrule" "test" {
  cpuload    = "1000"
  ramload    = "95"
  target     = "${data.abiquo_datacenter.test.id}"
}

data "abiquo_datacenter" "test" { name = "${var.datacenter}" }
