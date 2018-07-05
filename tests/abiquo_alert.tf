resource "abiquo_alert" "test" {
  virtualappliance = "${abiquo_vapp.test.id}"
  name        = "test alert"
  description = "test alert"

  alarms = [
    "${abiquo_alarm.test.id}"
  ]

  subscribers = [
    "test1@test.com",
    "test2@test.com",
  ]
}

data "abiquo_vdc"        "test"   { name = "tests" }
data "abiquo_template"   "test"   {
  templates = "${data.abiquo_vdc.test.templates}"
  name      = "tests"
}

resource "abiquo_vapp" "test" {
  virtualdatacenter = "${data.abiquo_vdc.test.id}"
  name              = "test alert"
}

resource "abiquo_vm" "test" {
  cpu                    = 1
  deploy                 = false
  ram                    = 64
  label                  = "test alert"
  virtualappliance       = "${abiquo_vapp.test.id}"
  virtualmachinetemplate = "${data.abiquo_template.test.id}"
}

resource "abiquo_alarm" "test" {
  target      = "${abiquo_vm.test.id}"
  name        = "test alert"
  metric      = "cpu_time"
  timerange   = 3
  statistic   = "average"
  formula     = "lessthan"
  threshold   = 10000
}
