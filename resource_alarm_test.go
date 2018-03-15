package main

import (
	"fmt"
	"testing"

	"github.com/abiquo/ojal/core"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

var testAccAlarmBasic = `
data "abiquo_enterprise" "enterprise" { name = "Abiquo" }
data "abiquo_location"   "location"   { name = "datacenter 1" }
data "abiquo_template"   "template"   { name = "tests" }

resource "abiquo_vdc" "vdc" {
	enterprise = "${data.abiquo_enterprise.enterprise.id}"
	location   = "${data.abiquo_location.location.id}"
	name       = "testAccPlanBasic"
	type       = "VMX_04"
}

resource "abiquo_vapp" "vapp" {
	virtualdatacenter = "${abiquo_vdc.vdc.id}"
	name              = "testAccPlanBasic"
}

resource "abiquo_vm" "vm" {
	cpu                    = 1
	deploy                 = false
	ram                    = 64
	label                  = "testAccPlanBasic"
	virtualappliance       = "${abiquo_vapp.vapp.id}"
	virtualmachinetemplate = "${data.abiquo_template.template.id}"
}

resource "abiquo_alarm" "alarm" {
  target      = "${abiquo_vm.vm.id}"
  name        = "testAccAlarmBasic"
  metric      = "cpu_time"
  period      = 60
  evaluations = 3
  statistic   = "average"
  formula     = "lessthan"
  threshold   = 10000
}
`

func TestAccAlarm_update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testCheckAlarmDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccAlarmBasic,
				Check: resource.ComposeTestCheckFunc(
					testCheckAlarmExists("abiquo_alarm.alarm"),
				),
			},
		},
	})
}

func testCheckAlarmDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "abiquo_alarm" {
			continue
		}
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "alarm")
		if err := core.Read(endpoint, nil); err == nil {
			return fmt.Errorf("alarm %q still exists", rs.Primary.Attributes["name"])
		}
	}
	return nil
}

func testCheckAlarmExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("alarm %q not found", name)
		}
		href := rs.Primary.Attributes["id"]
		endpoint := core.NewLinkType(href, "alarm")
		return core.Read(endpoint, nil)
	}
}
