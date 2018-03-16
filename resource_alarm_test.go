package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var alarmTestHelper = &testHelper{
	kind:  "abiquo_alarm",
	media: "alarm",
	config: `
	data "abiquo_enterprise" "test" { name = "Abiquo" }
	data "abiquo_location"   "test"   { name = "datacenter 1" }
	data "abiquo_template"   "test"   { name = "tests" }

	resource "abiquo_vdc" "test" {
		enterprise = "${data.abiquo_enterprise.test.id}"
		location   = "${data.abiquo_location.test.id}"
		name       = "testAccAlarmBasic"
		type       = "KVM"
	}

	resource "abiquo_vapp" "test" {
		virtualdatacenter = "${abiquo_vdc.test.id}"
		name              = "testAccAlarmBasic"
	}

	resource "abiquo_vm" "test" {
		cpu                    = 1
		deploy                 = false
		ram                    = 64
		label                  = "testAccAlarmBasic"
		virtualappliance       = "${abiquo_vapp.test.id}"
		virtualmachinetemplate = "${data.abiquo_template.test.id}"
	}

	resource "abiquo_alarm" "test" {
	  target      = "${abiquo_vm.test.id}"
	  name        = "testAccAlarmBasic"
	  metric      = "cpu_time"
	  period      = 60
	  evaluations = 3
	  statistic   = "average"
	  formula     = "lessthan"
	  threshold   = 10000
	}
	`,
}

func TestAccAlarm_update(t *testing.T) {
	resource.Test(t, alarmTestHelper.updateCase(t))
}
