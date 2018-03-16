package main

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

var alertTestHelper = &testHelper{
	kind:  "abiquo_alert",
	media: "alert",
	config: `
	data "abiquo_enterprise" "test" { name = "Abiquo" }
	data "abiquo_location"   "test"   { name = "datacenter 1" }
	data "abiquo_template"   "test"   { name = "tests" }

	resource "abiquo_vdc" "test" {
		enterprise = "${data.abiquo_enterprise.test.id}"
		location   = "${data.abiquo_location.test.id}"
		name       = "testAccAlertBasic"
		type       = "KVM"
	}

	resource "abiquo_vapp" "test" {
		virtualdatacenter = "${abiquo_vdc.test.id}"
		name              = "testAccAlertBasic"
	}

	resource "abiquo_vm" "test" {
		cpu                    = 1
		deploy                 = false
		ram                    = 64
		label                  = "testAccAlertBasic"
		virtualappliance       = "${abiquo_vapp.test.id}"
		virtualmachinetemplate = "${data.abiquo_template.test.id}"
	}

	resource "abiquo_alarm" "test" {
	  target      = "${abiquo_vm.test.id}"
	  name        = "testAccAlertBasic"
	  metric      = "cpu_time"
	  period      = 60
	  evaluations = 3
	  statistic   = "average"
	  formula     = "lessthan"
	  threshold   = 10000
	}

	resource "abiquo_alert" "test" {
	  virtualappliance = "${abiquo_vapp.test.id}"
	  name        = "testAccAlertBasic"
	  description = "testAccAlertBasic"
	  
	  alarms = [
	    "${abiquo_alarm.test.id}"
	  ]
	}
	`,
}

func TestAccAlert_update(t *testing.T) {
	resource.Test(t, alertTestHelper.updateCase(t))
}
