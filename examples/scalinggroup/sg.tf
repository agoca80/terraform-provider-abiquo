# Provider configuration: Get these from the environment
variable "endpoint" { }

provider "abiquo" {
  endpoint       = "${var.endpoint}"
  username       = "neutron"
  password       = "12qwaszx"
}

data "abiquo_location"   "location"   { name = "datacenter 1" }
data "abiquo_enterprise" "enterprise" { name = "neutron" }

resource "abiquo_vdc" "vdc" {
  enterprise = "${data.abiquo_enterprise.enterprise.id}"
  location   = "${data.abiquo_location.location.id}"

  name       = "tf example sg"
  type       = "KVM"
  
  cpusoft = 4 , cpuhard = 8
}

resource "abiquo_private" "private" {
  virtualdatacenter = "${abiquo_vdc.vdc.id}"

  # XXX workaround ABICLOUDPREMIUM-9660
  lifecycle = { ignore_changes = [ "dns1", "dns2" ] }

  mask    = 24
  address = "172.16.27.0"
  gateway = "172.16.27.1"
  name    = "terraform example-sg"
  dns1    = "8.8.8.8"
  dns2    = "4.4.4.4"
  suffix  = "lb.com"
}

resource "abiquo_fw" "fw" {
  virtualdatacenter = "${abiquo_vdc.vdc.id}"

  description = "terraform example-sg"
  name        = "terraform example-sg"

  # XXX workaround ABICLOUDPREMIUM-9668
  rules = [
    { protocol = "TCP", fromport = 22, toport = 22, sources = ["0.0.0.0/0"] },
    { protocol = "TCP", fromport = 80, toport = 80, sources = ["0.0.0.0/0"] },
    { protocol = "TCP", fromport = 443, toport = 443, sources = ["0.0.0.0/0"] }
  ]
}

resource "abiquo_lb" "lb" {
  virtualdatacenter = "${abiquo_vdc.vdc.id}"
  privatenetwork    = "${abiquo_private.private.id}"

  name         = "terraform"
  algorithm    = "ROUND_ROBIN"
  routingrules = [
    { protocolin = "HTTP" , protocolout = "HTTP" , portin = 80 , portout = 80 }
  ]
} 

resource "abiquo_vapp" "vapp" {
  virtualdatacenter = "${abiquo_vdc.vdc.id}"  
  name = "terraform"
}

data "abiquo_template" "template" { name = "tests" }

# Scaling group master instance
resource "abiquo_vm" "vm" {
  virtualappliance       = "${abiquo_vapp.vapp.id}"
  virtualmachinetemplate = "${data.abiquo_template.template.id}"

  label = "terraform"
  lbs = [ "${abiquo_lb.lb.id}" ]
  fws = [ "${abiquo_fw.fw.id}" ]  
}

resource "abiquo_alarm" "alarm0" {
  target      = "${abiquo_vm.vm.id}"
  name        = "terraform 0"
  metric      = "cpu_time"
  period      = 60
  evaluations = 3
  statistic   = "average"
  formula     = "lessthan"
  threshold   = 10000
}

resource "abiquo_alarm" "alarm1" {
  target      = "${abiquo_vm.vm.id}"
  name        = "terraform 1"
  metric      = "vcpu_time"
  period      = 60
  evaluations = 3
  statistic   = "average"
  formula     = "lessthan"
  threshold   = 10000
}

resource "abiquo_alert" "alert" {
  virtualappliance = "${abiquo_vapp.vapp.id}"
  name        = "terraform"
  description = "terraform"
  
  # XXX workaround entries order changes
  lifecycle = { ignore_changes = [ "alarms" ] }

  alarms = [
    "${abiquo_alarm.alarm0.id}",
    "${abiquo_alarm.alarm1.id}"
  ]
}

resource "abiquo_plan" "increase" {
  virtualmachine = "${abiquo_vm.vm.id}"
  name        = "increase"
  description = "increase"
  entries = [
    { parameter = "" , parametertype = "None", type = "SCALE_OUT" }
  ]
  
  triggers = [
    "${abiquo_alert.alert.id}"
  ]
}

resource "abiquo_plan" "decrease" {
  virtualmachine = "${abiquo_vm.vm.id}"
  name        = "decrease"
  description = "decrease"
  entries = [
    { parameter = "", parametertype = "None", type = "SCALE_IN" }
  ]
  
  triggers = [
    "${abiquo_alert.alert.id}"
  ]
}

resource "abiquo_sg" "sg" {
  mastervirtualmachine = "${abiquo_vm.vm.id}"
  virtualappliance     = "${abiquo_vapp.vapp.id}"

  name      = "terraform"
  cooldown  = 60
  min       = 2
  max       = 4
  scale_in  = [ { numberofinstances = 1 } ]
  scale_out = [ { numberofinstances = 1 } ]
  # scale_in  = [ { number = 1, from = "2017/01/01 00:00", until = "2019/01/01 00:00" } ]
  # scale_out = [ { number = 1, from = "2017/01/01 00:00", until = "2019/01/01 00:00" } ]
}

