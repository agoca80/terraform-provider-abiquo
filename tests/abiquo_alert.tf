resource "abiquo_alert" "test" {
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

data "abiquo_virtualdatacenter"        "test"   { name = "${var.virtualdatacenter}" }
data "abiquo_template"   "test"   {
  templates = "${data.abiquo_virtualdatacenter.test.templates}"
  name = "${var.template}"
}

resource "abiquo_virtualappliance" "test" {
  virtualdatacenter = "${data.abiquo_virtualdatacenter.test.id}"
  name              = "test alert"
}

resource "abiquo_virtualmachine" "test" {
  cpu                    = 1
  deploy                 = false
  ram                    = 64
  label                  = "test alert"
  virtualappliance       = "${abiquo_virtualappliance.test.id}"
  virtualmachinetemplate = "${data.abiquo_template.test.id}"
}

resource "abiquo_alarm" "test" {
  target      = "${abiquo_virtualmachine.test.id}"
  name        = "test alert"
  metric      = "${var.metric}"
  timerange   = 3
  statistic   = "average"
  formula     = "lessthan"
  threshold   = 10000
}
