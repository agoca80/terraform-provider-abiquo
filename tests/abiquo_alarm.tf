data "abiquo_template"   "test"   { name = "tests" }
data "abiquo_vdc"        "test"   { name = "tests" }

resource "abiquo_vapp" "test" {
  virtualdatacenter = "${data.abiquo_vdc.test.id}"
  name              = "test alarm"
}

resource "abiquo_vm" "test" {
  cpu                    = 1
  deploy                 = false
  ram                    = 64
  label                  = "test alarm"
  virtualappliance       = "${abiquo_vapp.test.id}"
  virtualmachinetemplate = "${data.abiquo_template.test.id}"
}

resource "abiquo_alarm" "test" {
  target      = "${abiquo_vm.test.id}"
  name        = "test alarm"
  metric      = "cpu_time"
  period      = 60
  evaluations = 3
  statistic   = "average"
  formula     = "lessthan"
  threshold   = 10000
}
